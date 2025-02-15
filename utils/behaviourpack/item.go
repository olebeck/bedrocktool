package behaviourpack

import "github.com/sandertv/gophertunnel/minecraft/protocol"

type itemDescription struct {
	Category       string `json:"category"`
	Identifier     string `json:"identifier"`
	IsExperimental bool   `json:"is_experimental"`
}

type minecraftItem struct {
	Description itemDescription `json:"description"`
	Components  map[string]any  `json:"components,omitempty"`
}

type itemBehaviour struct {
	FormatVersion string        `json:"format_version"`
	MinecraftItem minecraftItem `json:"minecraft:item"`
}

func (bp *BehaviourPack) AddItem(item protocol.ItemEntry) {
	ns, _ := ns_name_split(item.Name)
	if ns == "minecraft" {
		return
	}

	entry := itemBehaviour{
		FormatVersion: bp.formatVersion,
		MinecraftItem: minecraftItem{
			Description: itemDescription{
				Identifier:     item.Name,
				IsExperimental: true,
			},
			Components: make(map[string]any),
		},
	}
	bp.items[item.Name] = entry
}

func (bp *BehaviourPack) ApplyComponentEntries(entries []protocol.ItemComponentEntry) {
	for _, ice := range entries {
		item, ok := bp.items[ice.Name]
		if !ok {
			continue
		}
		item.MinecraftItem.Components = ice.Data
	}
}
