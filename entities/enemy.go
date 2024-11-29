package entities

import "rpg-go/components"

type Enemy struct {
	*Sprite
	FollowsPlayer bool
	CombatComp    *components.BasicCombat
}
