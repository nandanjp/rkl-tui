package main

type Pokemon struct {
	id   uint64
	name string
}

type Ability struct {
	id                  uint64
	name                string
	is_main_series      bool
	generation          Generation
	names               []Name
	effect_entries      []VerboseEffect
	effect_changes      []AbilityEffectChange
	flavor_text_entries []AbilityFlavorText
	pokemon             []AbilityPokemon
}

type AbilityEffectChange struct {
	effect_entries []Effect
	version_group  NamedApiResource
}

type AbilityFlavorText struct {
	flavor_text   string
	language      NamedApiResource
	version_group NamedApiResource
}

type AbilityPokemon struct {
	is_hidden bool
	slot      uint64
	pokemon   NamedApiResource
}

type Characteristic struct {
	id              uint64
	gene_modulo     uint64
	possible_values []uint64
	highest_stat    NamedApiResource
	descriptions    []Description
}

type EggGroup struct {
	id              uint64
	name            string
	names           []Name
	pokemon_species []NamedApiResource
}

type Gender struct {
	id                      uint64
	name                    string
	pokemon_species_details []PokemonSpeciesGender
	required_for_evolution  []NamedApiResource
}

type PokemonSpeciesGender struct {
	rate            uint64
	pokemon_species NamedApiResource
}

type VerboseEffect struct {
	effect       string
	short_effect string
	language     Language
}

type Language struct {
	id       uint64
	name     string
	official bool
	iso639   string
	iso3166  string
	names    []Name
}

type Name struct {
	name     string
	language Language
}

type Generation struct {
	id   uint64
	name string
}

type Move struct {
}
