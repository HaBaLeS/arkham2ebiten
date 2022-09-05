package renderer

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

func (ig *InvestigatorGui) Update() {

	//check if active -- aka you can click it

	//activePlayer := runtime.ScenarioSession().CurrentPlayer
	//activePlayer.Location -> Clues > 0

}

func (ig *InvestigatorGui) LoadGuiSprites() []*GuiSprite {

	ig.drawCard = NewGuiSprite("draw_card", "draw_card.png")
	ig.drawCard.X = 10
	ig.drawCard.Y = 1080 - 60
	ig.drawCard.OnClickFunc = ig.drawCard.onClickFuncDummy

	ig.engage = NewGuiSprite("engage", "engage.png")
	ig.engage.X = 10 + 1*205
	ig.engage.Y = 1080 - 60
	ig.engage.OnClickFunc = ig.engage.onClickFuncDummy

	ig.escape = NewGuiSprite("escape", "escape.png")
	ig.escape.X = 10 + 2*205
	ig.escape.Y = 1080 - 60
	ig.escape.OnClickFunc = ig.escape.onClickFuncDummy

	ig.fight = NewGuiSprite("fight", "fight.png")
	ig.fight.X = 10 + 3*205
	ig.fight.Y = 1080 - 60
	ig.fight.OnClickFunc = ig.fight.onClickFuncDummy

	ig.investigate = NewGuiSprite("investigate", "investigate.png")
	ig.investigate.X = 10 + 4*205
	ig.investigate.Y = 1080 - 60
	ig.investigate.OnClickFunc = ig.investigate.onClickFuncDummy

	ig.move = NewGuiSprite("move", "move.png")
	ig.move.X = 10 + 5*205
	ig.move.Y = 1080 - 60
	ig.move.OnClickFunc = ig.move.onClickFuncDummy

	ig.playCard = NewGuiSprite("play_card", "play_card.png")
	ig.playCard.X = 10 + 6*205
	ig.playCard.Y = 1080 - 60
	ig.playCard.OnClickFunc = ig.playCard.onClickFuncDummy

	ig.resource = NewGuiSprite("resource", "resource.png")
	ig.resource.X = 10 + 7*205
	ig.resource.Y = 1080 - 60
	ig.resource.OnClickFunc = ig.resource.onClickFuncDummy

	retList := make([]*GuiSprite, 0)
	retList = append(retList, ig.escape, ig.drawCard, ig.engage, ig.fight, ig.move, ig.playCard, ig.resource, ig.investigate)

	return retList
}

func (ig *InvestigatorGui) Enable() {
	ig.drawCard.Visible()
	ig.drawCard.Active()

	ig.engage.Visible()
	ig.engage.Inactive()
	ig.escape.Visible()
	ig.escape.Inactive()
	ig.fight.Visible()
	ig.fight.Inactive()
	ig.investigate.Visible()
	ig.investigate.Visible()
	ig.move.Visible()
	ig.move.Inactive()
	ig.playCard.Visible()
	ig.playCard.Active()
	ig.resource.Visible()
	ig.resource.Active()
}
