/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package operator

// Operator condition type
type Operator string

const (
	// Tr Tree condition
	Tr Operator = "tr"

	// Query Selectors

	// ========== Comparison =============

	// Eq Matches values that are equal to a specified value
	Eq Operator = "eq"
	// Ne Matches all values that are not equal to a specified value
	Ne Operator = "ne"
	// Lt Matches values that are less than a specified value
	Lt Operator = "lt"
	// Lte Matches values that are less than or equal to a specified value
	Lte Operator = "lte"
	// Gt Matches values that are greater than a specified value
	Gt Operator = "gt"
	// Gte Matches values that are greater than or equal to a specified value
	Gte Operator = "gte"
	// In Matches any of the values specified in an array
	In Operator = "in"
	// Nin Matches none of the values specified in an array
	Nin Operator = "nin"
	// Con Matches values that contains some string
	Con Operator = "con"

	// ========== Logical =============

	// Not Inverts the effect of a query expression and returns documents that do not match the query expression
	Not Operator = "not"
	// And Joins query clauses with a logical AND returns all documents that match the conditions of both clauses
	And Operator = "and"
	// Or Joins query clauses with a logical OR returns all documents that match the conditions of either clause
	Or Operator = "or"
	// Nor Joins query clauses with a logical NOR returns all documents that fail to match both clauses
	Nor Operator = "nor"

	// ========== Element =============

	// Ext Matches documents that have the specified field
	Ext Operator = "exists"
	// Typ Selects documents if a field is of the specified type
	Typ Operator = "type"

	// ========== Pipeline in ChangeStream =========

	// Mat matches documents in change stream
	Mat Operator = "match"
)

var (
	// EmptyCondition just return {}
	EmptyCondition = NewLeafCondition(Tr, M{})
)

// M struct to store map data
type M map[string]interface{}

// Update update M
func (m M) Update(key string, value interface{}) M {
	m[key] = value
	return m
}

// Merge merge additional key-value pair into M
func (m M) Merge(additionM M) {
	for key, value := range additionM {
		m[key] = value
	}
}

// Condition the condition making is just the process of tree building
type Condition struct {
	Op       Operator
	Value    interface{}
	Children []*Condition
}

// NewLeafCondition create leaf condition
func NewLeafCondition(t Operator, v interface{}) *Condition {
	return &Condition{
		Op:    t,
		Value: v,
	}
}

// NewBranchCondition create branch condition
func NewBranchCondition(t Operator, cons ...*Condition) *Condition {
	newCondition := &Condition{
		Op: t,
	}
	newCondition.Children = append(newCondition.Children, cons...)
	return newCondition
}

// Combine combine
func (c *Condition) Combine(leafFunc func(Operator, interface{}) interface{},
	combineFunc func(Operator, []*Condition) interface{}) interface{} {
	// leaf node
	if len(c.Children) == 0 {
		return leafFunc(c.Op, c.Value)
	}
	return combineFunc(c.Op, c.Children)
}

// import (
// 	// "container/list"
// 	// "reflect"
// )

// // Operator condition type
// type Operator string

// const (
// 	// Tr Tree condition
// 	Tr Operator = "tr"

// 	// Query Selectors

// 	// ========== Comparison =============

// 	// Eq Matches values that are equal to a specified value
// 	Eq Operator = "eq"
// 	// Ne Matches all values that are not equal to a specified value
// 	Ne Operator = "ne"
// 	// Lt Matches values that are less than a specified value
// 	Lt Operator = "lt"
// 	// Lte Matches values that are less than or equal to a specified value
// 	Lte Operator = "lte"
// 	// Gt Matches values that are greater than a specified value
// 	Gt Operator = "gt"
// 	// Gte Matches values that are greater than or equal to a specified value
// 	Gte Operator = "gte"
// 	// In Matches any of the values specified in an array
// 	In Operator = "in"
// 	// Nin Matches none of the values specified in an array
// 	Nin Operator = "nin"
// 	// Con Matches values that contains some string
// 	Con Operator = "con"

// 	// ========== Logical =============

// 	// Not Inverts the effect of a query expression and returns documents that do not match the query expression
// 	Not Operator = "not"
// 	// And Joins query clauses with a logical AND returns all documents that match the conditions of both clauses
// 	And Operator = "and"
// 	// Or Joins query clauses with a logical OR returns all documents that match the conditions of either clause
// 	Or Operator = "or"
// 	// Nor Joins query clauses with a logical NOR returns all documents that fail to match both clauses
// 	Nor Operator = "nor"

// 	// ========== Element =============

// 	// Ext Matches documents that have the specified field
// 	Ext Operator = "exists"
// 	// Typ Selects documents if a field is of the specified type
// 	Typ Operator = "type"

// 	// ========== Pipeline in ChangeStream =========

// 	// Mat matches documents in change stream
// 	Mat Operator = "match"
// )

// // Condition describe the conditions of database operations
// // It is a chain-tree structure: basically a tree but every node is a chain.
// // There is two types of Condition according to Condition.Value:
// //
// //  - leaf node:   whose value is operator.M, the actual condition
// //  - branch node: whose value contains 2 types:
// //         - Condition:     "Not" operation build it, has 1 child node
// //         - ConditionPair: "And/Or" operation build it, has 2 child node
// //
// // Let Cm = leaf node, Cc = branch node(Condition), Cp = branch node(Condition Pair)
// // See the following struct:
// //
// //    ---------------------
// //    | Cm1 -> Cp1 -> Cm2 |
// //    ----------|----------
// //              |(and/or)
// //     _________|___________
// //     |                   |
// // --------------   ---------------------
// // | Cm3 -> Cm4 |   | Cc1 -> Cm5 -> Cc2 |
// // --------------   ---|-------------|---
// //                     |(not)        |(not)
// //                  -------        --------------
// //                  | Cm6 |        | Cm7 -> Cm8 |
// //                  -------        --------------
// //
// // Every "block" contains 1 "chain".
// // Every "chain" contains some condition "node"
// // Every "node": contains 3 types:
// //     - Cm(leaf node):   actual condition
// //     - Cp(branch node): and/or operation, link to 2 "block"
// //     - Cc(branch node): not operation, link to 1 "block"
// //
// // The condition making is just the process of tree building
// type Condition struct {
// 	Op       Operator
// 	Value    interface{}
// 	Children []*Condition
// }

// // ConditionPair pair of condition
// type ConditionPair struct {
// 	First  *Condition
// 	Second *Condition
// }

// var (
// // BaseCondition is a True node of condition chain
// // and can be called like: operator.BaseCondition.AddOp(cType, key, value)
// // then the func will return a new Condition which contains a condition chain
// // with a BaseCondition head.
// //  BaseCondition -> Condition1 -> Condition2
// //                -> Condition3
// //                -> Condition4
// // See the Condition.AddOp() for details
// // ATTENTION: BaseCondition is out of condition tree, except the tree is empty(then will be always true).
// //BaseCondition = NewCondition(Tr, M(nil))
// )

// // NewLeafCondition create leaf condition
// func NewLeafCondition(t Operator, v interface{}) *Condition {
// 	return &Condition{
// 		Op:    t,
// 		Value: v,
// 	}
// }

// // NewBranchCondition create branch condition
// func NewBranchCondition(t Operator, cons ...*Condition) *Condition {
// 	newCondition := &Condition{
// 		Op: t,
// 	}
// 	newCondition.Children = append(newCondition.Children, cons...)
// 	return newCondition
// }

// NewCondition create new condition
// func NewCondition(t ConditionType, v interface{}) (r *Condition) {
// 	r = &Condition{
// 		list:  list.New(),
// 		Type:  t,
// 		Value: v,
// 	}
// 	r.list.PushBack(r)
// 	return
// }

// And return a new Condition which is made up of two Condition with AND
// if either of then is BaseCondition then return another
// func (c *Condition) And(s *Condition) *Condition {
// 	if reflect.DeepEqual(c, BaseCondition) {
// 		return s
// 	}
// 	if reflect.DeepEqual(s, BaseCondition) {
// 		return c
// 	}
// 	return NewCondition(And, &ConditionPair{c, s})
// }

// Or return a new Condition which is made up of two Condition with OR
// if either of then is BaseCondition then return another
// func (c *Condition) Or(s *Condition) *Condition {
// 	if reflect.DeepEqual(c, BaseCondition) {
// 		return s
// 	}
// 	if reflect.DeepEqual(s, BaseCondition) {
// 		return c
// 	}
// 	return NewCondition(Or, &ConditionPair{c, s})
// }

// Not return a new Condition which is made up of itself with NOT
// if it is BaseCondition then return itself
// Not(BaseCondition) -> BaseCondition
// func (c *Condition) Not() *Condition {
// 	if reflect.DeepEqual(c, BaseCondition) {
// 		return c
// 	}
// 	return NewCondition(Not, c)
// }

// AddOp return a new Condition linked to the old one, just like append the chain with a new node
// func (c *Condition) AddOp(t ConditionType, key string, value interface{}) (r *Condition) {
// 	r = NewCondition(t, M{key: value})
// 	if reflect.DeepEqual(c, BaseCondition) {
// 		return
// 	}
// 	r.list.PushFrontList(c.list)
// 	return
// }

// Combine provide a common process that combine the Condition tree.
// Any driver implements the Tank can combine the Condition tree by providing two functions:
//
//  - leafFunc: handle the leaf node whose value-type is M and it's the actual condition
//  - combineFunc: handle the branch node whose value-type is Condition or ConditionPair
//
// Combine should be called before doing any filter-need operation like query, update, remove. And it should be cached by
// driver to prevent processing every time.
// func (c *Condition) Combine(leafFunc func(Operator, interface{}) interface{}, combineFunc func(Operator, []*Condition) interface{}) interface{} {
// 	// leaf node
// 	if len(c.Children) == 0 {
// 		return leafFunc(c.Operator, c.Value)
// 	}
// 	return combineFunc(c.Operator, c.Children)
//
// if c.list == nil {
// 	return nil
// }

// tv := make([]interface{}, 0, c.list.Len())
// for e := c.list.Front(); e != nil; e = e.Next() {
// 	branch := e.Value.(*Condition)

// 	var tmp interface{}
// 	switch branch.Value.(type) {

// 	// Leaf node condition combined with AddOp() or NewCondition()
// 	case M:
// 		tmp = leafFunc(branch)

// 		// Single condition combined with Not()
// 	case *Condition:
// 		child := branch.Value.(*Condition)
// 		tmp = child.Combine(leafFunc, combineFunc)
// 		//tmp = combineFunc(Not, []interface{}{tmpNot})
// 		tmp = combineFunc{}

// 		// Condition pair that combined with And() or Or()
// 	case *ConditionPair:
// 		child := branch.Value.(*ConditionPair)
// 		tmpFirst := child.First.Combine(leafFunc, combineFunc)
// 		tmpSecond := child.Second.Combine(leafFunc, combineFunc)
// 		tmp = combineFunc(branch.Type, []interface{}{tmpFirst, tmpSecond})
// 	}

// 	tv = append(tv, tmp)
// }
// return combineFunc(And, tv)
// }
