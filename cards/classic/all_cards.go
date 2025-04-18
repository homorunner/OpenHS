package cards

var BasicHeros = []interface{}{
	&Jaina{},
	&Rexxar{},
	&Anduin{},
	&Garrosh{},
	&Thrall{},
	&Uther{},
	&Valeera{},
	&Malfurion{},
	&Guldan{},
}

var AllCards = append(BasicHeros, []interface{}{
	&WaterElemental{},
}...)
