package main

import (
	"fmt"
	"github.com/janmbaco/go-redux/src"
)

type Actions struct {
	Sum          redux.Action
	Substraction redux.Action
}

func Sum(state int, payload int) int {
	return state + payload
}

type SubstractionLogic struct {
}

func (r *SubstractionLogic) Substraction(state int, payload int) int {
	return state - payload
}

func main() {

	actions := &Actions{}

	builder := redux.NewBusinessObjectBuilder(0, actions)
	builder.On(actions.Sum, Sum)
	builder.SetActionsLogicByObject(&SubstractionLogic{})
	builder.SetSelector("counter")

	counterBO := builder.GetBusinessObject()

	store := redux.NewStore(counterBO)

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: '0'

	store.Dispatch(actions.Sum.With(1))
	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: '1'

	globalSubscription := func() {
		fmt.Printf("globalSubscription - state changed, current state: '%v'\n", store.GetState())
	}
	store.Subscribe(&globalSubscription)
	store.Dispatch(actions.Substraction.With(1))
	// output:
	// globalSubscription - state changed, current state: '0'

	store.UnSubscribe(&globalSubscription)

	actions2 := &Actions{}
	counter2BO := redux.NewBusinessObjectBuilder(10, actions2).
		On(actions2.Sum, Sum).
		SetActionsLogicByObject(&SubstractionLogic{}).
		SetSelector("counter2").
		GetBusinessObject()

	store.AddBusinessObject(counter2BO)

	actions3 := &Actions{}
	counter3BO := redux.NewBusinessObjectBuilder(100, actions3).
		On(actions3.Sum, Sum).
		SetActionsLogicByObject(&SubstractionLogic{}).
		SetSelector("counter3").
		GetBusinessObject()

	store.AddBusinessObject(counter3BO)

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	// current state: '0'

	store.Dispatch(actions.Sum.With(1))
	store.Dispatch(actions2.Sum.With(1))
	store.Dispatch(actions3.Sum.With(1))

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: '[1 11 101]'

	store.RemoveBusinessObject(counter3BO)
	func() {
		defer func() {
			if re := recover(); re != nil {
				fmt.Printf(re.(string) + "\n")
			}
		}()
		store.Dispatch(actions3.Sum.With(1))
	}()
	// output:
	//There is not any Reducers that execute this action!

	fmt.Printf("current state: '%v'\n", store.GetState())
	// output:
	//current state: '[1 11 101]'

	fmt.Printf("current state counter: '%v'\n", store.GetStateOf("counter"))
	fmt.Printf("current state counter2: '%v'\n", store.GetStateOf("counter2"))
	fmt.Printf("current state counter3: '%v'\n", store.GetStateOf("counter3"))
	// output:
	//current state counter: '1'
	//current state counter2: '11'
	//current state counter3: '101'

	counterSubscribe := func(newState interface{}) {
		fmt.Printf("counterSubscribe - state changed, current state: '%v'\n", newState)
	}
	store.SubscribeTo("counter", &counterSubscribe)
	store.Dispatch(actions.Substraction.With(1))
	store.Dispatch(actions2.Substraction.With(11))
	// output:
	//counterSubscribe - state changed, current state: '0'

	counter2Subscribe := func(newState interface{}) {
		fmt.Printf("counter2Subscribe - state changed, current state: '%v'\n", newState)
	}
	store.SubscribeTo("counter2", &counter2Subscribe)
	store.Dispatch(actions.Sum.With(1))
	store.Dispatch(actions2.Sum.With(10))
	// output:
	//ounterSubscribe - state changed, current state: '1'
	//counter2Subscribe - state changed, current state: '10'

	store.UnSubscribeFrom("counter", &counterSubscribe)
	store.Dispatch(actions.Sum.With(5))
	store.Dispatch(actions2.Substraction.With(8))
	// output:
	//counter2Subscribe - state changed, current state: '2'

	fmt.Printf("current state counter: '%v'\n", store.GetStateOf("counter"))
	// output:
	//current state counter: '6'

}
