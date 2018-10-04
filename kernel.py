import math
import numba
import sqlite3
import numpy as np
#import functools
#for nXn=>mXm where n>m
patience = 0.8
FIRST_HAND = 1
LAST_HAND = -1
SELF=1
ENEMY=-1
INIT_MAX_UTIL = -100
INIT_MIN_UTIL = 100
INIT_POLICY = -100

def pick_slice(target,index,sliceSize):
    result = []
    shift = 0
    targetSize = int(math.sqrt(len(target)))
    for _ in range(sliceSize):
        result+=target[index+shift:index+shift+sliceSize]
        shift+=targetSize
    return result

def num_of_slice_in_1D(target,slice):
    targetSize = int(math.sqrt(len(target)))
    return targetSize - slice +1

def size_of_target(target):
    return int(math.sqrt(len(target)))

def is_target_empty(target):
    return sum([0 if elem == 0 else 1 for elem in target]) == 0

def is_target_full(target):
    return (0 in target)

def check_row(target):
    step = size_of_target(target)
    temp = [sum(target[i: i+step]) for i in range(0,len(target),step)]
    result = (False,None)
    if (size_of_target(target) in temp):
        result = (False,1)
    elif (-size_of_target(target) in temp):
        result = (False,-1)
    elif (0 in target):
        result = (True,0)
    else:
        result = (False,0)
    return result

def check_column(target):
    step = size_of_target(target)
    [temp,*rest] = [target[i: i+step] for i in range(0,len(target),step)]
    for i in range(len(temp)):
        for line in rest:
            temp[i]+=line[i]
    result = (False,0)
    if (size_of_target(target) in temp):
        result = (False,1)
    elif (-size_of_target(target) in temp):
        result = (False,-1)
    elif (0 in target):
        result = (True,0)
    else:
        result = (False,0)
    return result

def check_diagnal(target):
    size = size_of_target(target)
    temp = sum([target[i*size_of_target(target)+i] for i in range(size)])
    result = (False,0)
    if temp == size:
        result = (False,1)
    elif temp == -size:
        result = (False,-1)
    elif 0 in target:
        result = (True, 0)
    else:
        result = (False,0)
    return result

def check_rev_diagnal(target):
    size = size_of_target(target)
    temp = sum([target[i*size_of_target(target)+(size - i - 1)] for i in range(size)])
    result = (False,0)
    if temp == size:
        result = (False,1)
    elif temp == -size:
        result = (False,-1)
    elif 0 in target:
        result = (True, 0)
    else:
        result = (False,0)
    return result

def terminal_state_check(target):
    (is_playable,winner) = check_row(target)
    if winner == 0 and is_playable:
        (is_playable,winner) = check_column(target)
    elif winner == 0 and is_playable:
        (is_playable, winner) = check_diagnal(target)
    elif winner == 0 and is_playable:
        (is_playable, winner) = check_rev_diagnal(target)
    return (is_playable,winner)
#r=0
def is_action_avalaible(state,action):
    return state[action] == 0

def apply_policy(state,policy,label):
    result = [s for s in state]
    if is_action_avalaible(state,policy):
        result[policy] = label
    return result

def generate_policy(target):
    return [policy for policy in range(len(target)) if target[policy] == 0]

#@numba.jit(nopython=True)
def best_policy_and_util(board, label,depth = 1,alpha = -math.inf,beta = math.inf,synmmatric_optimization = None):
    possible_policies = generate_policy(board)
    best_policy = INIT_POLICY
    #print(depth,"invoked best")
    state_pool = 
    (has_empty, winner) = terminal_state_check(board)
    max_util = INIT_MAX_UTIL
    if (has_empty and winner == 0):

        for policy in possible_policies:
            if (alpha<beta):
                image_board = apply_policy(board, policy, label)
                (_, r) = worst_policy_and_util(image_board, -label, depth + 1,alpha,beta)
                alpha = r if alpha<r else alpha
                if (r > max_util):
                    best_policy = policy
                    max_util = r
    else:
        max_util = winner * 10
    return (best_policy, max_util)


"""
    Min method
"""

#@numba.jit(nopython=True)
def worst_policy_and_util(board, label,depth = 1,alpha = -math.inf,beta = math.inf,synmmatric_optimization = None):
    possible_policies = generate_policy(board)
    worst_policy = INIT_POLICY
    #print(depth,"invoked worst")
    (has_empty, winner) = terminal_state_check(board)
    min_util = INIT_MIN_UTIL
    if (has_empty and winner == 0):

        for policy in possible_policies:
            if (alpha<beta):
                image_board = apply_policy(board, policy, label)
                (_, temp_util) = best_policy_and_util(image_board, -label, depth + 1,alpha,beta)
                beta = temp_util if temp_util<beta else beta
                if (temp_util < min_util):
                    worst_policy = policy
                    min_util = temp_util

    else:
        min_util = winner * 10

    return (worst_policy, min_util)


def policy_generalization(size_of_source,size_of_target,global_index,local_index):
    return (int(local_index/size_of_target)+int(global_index/size_of_source))*size_of_source+(local_index%size_of_target+global_index%size_of_source)
def policy2D(index,size):
    return (int(index/size),index%size)
def policy1D(policy,size):
    (dim1,dim2) = policy
    return dim1*size+dim2
def policy_rotate_90_deg(policy,size):#90deg
    (dim2,dim1) = policy2D(policy,size)
    return dim1*size+dim2
def policy_rotate_180_deg(policy,size):#90deg
    (dim1,dim2) = policy2D(policy,size)
    return (dim1)*size+size-dim2-1
def policy_rotate_270_deg(policy,size):#90deg
    (dim2,dim1) = policy2D(policy,size)
    return (size-dim1-1)*size+dim2

def symatric_state_inclusive(state):
    size = int(len(state)**0.5)
    return {0:state,
            90:list(np.rot90(np.reshape(state,(size,size)),k=1).flatten()),
            180:list(np.rot90(np.reshape(state,(size,size)),k=2).flatten()),
            270:list(np.rot90(np.reshape(state,(size,size)),k=3).flatten())}
#print (best_policy_and_util([1,1,1,0,0,-1,0,0,-1,0,0,0,-1,0,0,-1,0,0,0,0,0,0,0,0,0],1))
#print (policy_generalization(5,3,12,4))
#print (policy_rotate_180_deg(6,5))
print (symatric_state_inclusive([1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4]))



