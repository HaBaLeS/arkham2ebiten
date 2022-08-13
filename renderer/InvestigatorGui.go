package renderer

import "log"

type InvestigatorGui struct {
	drawCard    *GuiSprite
	engage      *GuiSprite
	escape      *GuiSprite
	fight       *GuiSprite
	investigate *GuiSprite
	move        *GuiSprite
	playCard    *GuiSprite
	resource    *GuiSprite
}

func (ig *InvestigatorGui) LoadGuiSprites() []*GuiSprite {

	ig.drawCard = NewGuiSprite("draw_card", "draw_card.png")
	ig.drawCard.X = 10
	ig.drawCard.Y = 1080 - 60
	ig.drawCard.OnClickFunc = onClickFuncDummy

	ig.engage = NewGuiSprite("engage", "engage.png")
	ig.engage.X = 10 + 1*205
	ig.engage.Y = 1080 - 60
	ig.engage.OnClickFunc = onClickFuncDummy

	ig.escape = NewGuiSprite("escape", "escape.png")
	ig.escape.X = 10 + 2*205
	ig.escape.Y = 1080 - 60
	ig.escape.OnClickFunc = onClickFuncDummy

	ig.fight = NewGuiSprite("fight", "fight.png")
	ig.fight.X = 10 + 3*205
	ig.fight.Y = 1080 - 60
	ig.fight.OnClickFunc = onClickFuncDummy

	ig.investigate = NewGuiSprite("investigate", "investigate.png")
	ig.investigate.X = 10 + 4*205
	ig.investigate.Y = 1080 - 60
	ig.investigate.OnClickFunc = onClickFuncDummy

	ig.move = NewGuiSprite("move", "move.png")
	ig.move.X = 10 + 5*205
	ig.move.Y = 1080 - 60
	ig.move.OnClickFunc = onClickFuncDummy

	ig.playCard = NewGuiSprite("play_card", "play_card.png")
	ig.playCard.X = 10 + 6*205
	ig.playCard.Y = 1080 - 60
	ig.playCard.OnClickFunc = onClickFuncDummy

	ig.resource = NewGuiSprite("resource", "resource.png")
	ig.resource.X = 10 + 7*205
	ig.resource.Y = 1080 - 60
	ig.resource.OnClickFunc = onClickFuncDummy

	retList := make([]*GuiSprite, 0)
	retList = append(retList, ig.escape, ig.drawCard, ig.engage, ig.fight, ig.move, ig.playCard, ig.resource, ig.investigate)

	return retList
}

func (ig *InvestigatorGui) Enable() {
	ig.drawCard.Enable()
	ig.engage.Enable()
	ig.escape.Enable()
	ig.fight.Enable()
	ig.investigate.Enable()
	ig.move.Enable()
	ig.playCard.Enable()
	ig.resource.Enable()
}

func onClickFuncDummy() {
	//fixme send command -> Gui which action was pressed
	//Engine should reduce number of possible actions
	//engine should send a disable gui after the numer == 0
		
	log.Printf("execute some stuff")
}
