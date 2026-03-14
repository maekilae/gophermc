// Code generated from Pkl module `makila.minecraftgo.properties`. DO NOT EDIT.
package config

// The Minecraft server configuration template.
type ServerProperties struct {
	// Allows users to use flight on the server while in Survival mode.
	AllowFlight bool `pkl:"AllowFlight"`

	// Allows players to travel to the Nether.
	AllowNether bool `pkl:"AllowNether"`

	// Send console command outputs to all online operators.
	BroadcastConsoleToOps bool `pkl:"BroadcastConsoleToOps"`

	// Send rcon console command outputs to all online operators.
	BroadcastRconToOps bool `pkl:"BroadcastRconToOps"`

	// Options: peaceful, easy, normal, hard
	Difficulty string `pkl:"Difficulty"`

	// Enables command blocks.
	EnableCommandBlock bool `pkl:"EnableCommandBlock"`

	// Exposes an MBean for JMX monitoring.
	EnableJmxMonitoring bool `pkl:"EnableJmxMonitoring"`

	// Enables remote access to the server console.
	EnableRcon bool `pkl:"EnableRcon"`

	// Makes the server appear as "online" on the server list.
	EnableStatus bool `pkl:"EnableStatus"`

	// Enables GameSpy4 protocol server listener.
	EnableQuery bool `pkl:"EnableQuery"`

	// Requires Mojang-signed public keys for connection.
	EnforceSecureProfile bool `pkl:"EnforceSecureProfile"`

	// Enforces the whitelist on the server.
	EnforceWhitelist bool `pkl:"EnforceWhitelist"`

	// Distance for entity rendering (10-1000).
	EntityBroadcastRangePercentage int `pkl:"EntityBroadcastRangePercentage"`

	// Force players to join in the default game mode.
	ForceGamemode bool `pkl:"ForceGamemode"`

	// Permission level for functions (1-4).
	FunctionPermissionLevel int `pkl:"FunctionPermissionLevel"`

	// Defines the mode of gameplay.
	Gamemode string `pkl:"Gamemode"`

	// Defines whether structures can be generated.
	GenerateStructures bool `pkl:"GenerateStructures"`

	// Settings used to customize world generation.
	GeneratorSettings string `pkl:"GeneratorSettings"`

	// Hardcore mode (spectator on death).
	Hardcore bool `pkl:"Hardcore"`

	// Hide player list on status requests.
	HideOnlinePlayers bool `pkl:"HideOnlinePlayers"`

	// Datapacks to disable/enable on creation.
	InitialDisabledPacks string `pkl:"InitialDisabledPacks"`

	InitialEnabledPacks string `pkl:"InitialEnabledPacks"`

	// World name and folder name.
	LevelName string `pkl:"LevelName"`

	LevelSeed string `pkl:"LevelSeed"`

	LevelType string `pkl:"LevelType"`

	// Max consecutive neighbor updates.
	MaxChainedNeighborUpdates int `pkl:"MaxChainedNeighborUpdates"`

	// Maximum simultaneous players.
	MaxPlayers int `pkl:"MaxPlayers"`

	// Watchdog timeout in milliseconds.
	MaxTickTime int `pkl:"MaxTickTime"`

	// Maximum world border radius.
	MaxWorldSize int `pkl:"MaxWorldSize"`

	// Message of the Day.
	Motd string `pkl:"Motd"`

	// Packet compression threshold.
	NetworkCompressionThreshold int `pkl:"NetworkCompressionThreshold"`

	// Authenticate players via Mojang.
	OnlineMode bool `pkl:"OnlineMode"`

	// Default permission level for operators.
	OpPermissionLevel int `pkl:"OpPermissionLevel"`

	// Kick idle players after X minutes.
	PlayerIdleTimeout int `pkl:"PlayerIdleTimeout"`

	// Prevent proxy/VPN connections.
	PreventProxyConnections bool `pkl:"PreventProxyConnections"`

	// Enable chat preview features.
	PreviewsChat bool `pkl:"PreviewsChat"`

	// Enable Player vs Player.
	Pvp bool `pkl:"Pvp"`

	// Ports for Query and RCON.
	QueryPort int `pkl:"QueryPort"`

	RconPassword string `pkl:"RconPassword"`

	RconPort int `pkl:"RconPort"`

	// Resource pack settings.
	ResourcePack string `pkl:"ResourcePack"`

	ResourcePackPrompt string `pkl:"ResourcePackPrompt"`

	ResourcePackSha1 string `pkl:"ResourcePackSha1"`

	RequireResourcePack bool `pkl:"RequireResourcePack"`

	// Network binding settings.
	ServerIp string `pkl:"ServerIp"`

	ServerPort int `pkl:"ServerPort"`

	// Chunk distance for updates and viewing.
	SimulationDistance int `pkl:"SimulationDistance"`

	ViewDistance int `pkl:"ViewDistance"`

	// Spawning toggles.
	SpawnAnimals bool `pkl:"SpawnAnimals"`

	SpawnMonsters bool `pkl:"SpawnMonsters"`

	SpawnNpcs bool `pkl:"SpawnNpcs"`

	SpawnProtection int `pkl:"SpawnProtection"`

	// Use native Linux transport for performance.
	UseNativeTransport bool `pkl:"UseNativeTransport"`

	// Enables a whitelist.
	Whitelist bool `pkl:"Whitelist"`
}
