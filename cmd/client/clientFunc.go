package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func HandlerPause(gameState *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		defer fmt.Printf("> ")
		gameState.HandlePause(ps)
		return pubsub.Ack
	}
}

func HandlerMove(gameState *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Printf("> ")
		resultMove := gameState.HandleMove(move)
		switch resultMove {
		case gamelogic.MoveOutComeSafe:
			fallthrough
		case gamelogic.MoveOutcomeMakeWar:
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			fallthrough
		default:
			return pubsub.NackDiscard
		}
	}
}
