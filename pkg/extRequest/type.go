package extRequest

type rdvStc struct {
	Page pageStc `xml:"PAGE"`
}

type pageStc struct {
	Lock statut `xml:"LOCK"`
	Session statut `xml:"SESSION"`
	Metadata metadata `xml:"METADATA"`
	Cg value `xml:"CG"`
	Ccg value `xml:"CCG"`
	Poste poste `xml:"POSTE"`
}

type statut struct {
	Statut string `xml:"STATUT"`
}

type metadata struct {
	Lang string `xml:"LANG"`
	IdPage int64 `xml:"ID_PAGE"`
}

type value struct {
	Value string `xml:"VALUE"`
}

type poste struct {
	Ferm value `xml:"FERM"`
}
