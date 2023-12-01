package plugins

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/joho/godotenv"
)

type PluginManager struct {
	pluginMap map[string]IPlugin
	methodMap map[string]IPlugin
}

func (p *PluginManager) InitializePlugins() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var pluginsDirectory = os.Getenv("PLUGINS_DIRECTORY")

	p.loadPlugins(pluginsDirectory)
	p.initializePlugins()
}

func (p *PluginManager) CallMethod(methodName string, params ...interface{}) (interface{}, error) {
	plugin := p.methodMap[methodName]
	if plugin == nil {
		return nil, fmt.Errorf("Method not registered in any of plugins %v", methodName)
	}
	return plugin.CallMethod(methodName, params...)
}

func (p *PluginManager) loadPlugins(pluginsDirectory string) {
	p.pluginMap = make(map[string]IPlugin)

	plugins, err := getPluginFiles(pluginsDirectory, ".so")
	if err != nil {
		fmt.Println("error when searching for plugins")
	}

	for _, pluginPath := range plugins {
		pluginName := strings.TrimSuffix(filepath.Base(pluginPath), filepath.Ext(pluginPath))

		plugin, err := loadPlugin(pluginPath)
		if err != nil {
			fmt.Println("Error loading plugin ", pluginPath, err)
		}

		p.pluginMap[pluginName] = plugin
	}
}

func (p *PluginManager) initializePlugins() {
	p.methodMap = make(map[string]IPlugin)

	for pluginName, plugin := range p.pluginMap {
		err := plugin.Initialize()
		if err != nil {
			fmt.Println("Error initializing plugin ", pluginName, err)
		}

		methods, err := plugin.GetMethods()
		if err != nil {
			fmt.Println("Error getting methods from plugin ", pluginName, err)
		}

		// Register method for plugin
		for _, method := range methods {
			if p.methodMap[method.MethodName] != nil {
				fmt.Println("Duplicate method found", method.MethodName, p.methodMap[method.MethodName], plugin)
				continue
			}
			p.methodMap[method.MethodName] = plugin
		}
	}
}

func getPluginFiles(pluginsDirectory, extension string) ([]string, error) {
	var result []string

	err := filepath.WalkDir(pluginsDirectory, func(fileName string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		// Skip directory
		if dirEntry.IsDir() {
			return nil
		}

		// Check for expected extension
		if strings.HasSuffix(dirEntry.Name(), extension) {
			result = append(result, fileName)
			return nil
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error searching plugins in directory %v: %v", pluginsDirectory, err)
	}

	return result, nil
}

func loadPlugin(path string) (IPlugin, error) {
	plugin, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	sym, err := plugin.Lookup("ExportPlugin")
	if err != nil {
		return nil, err
	}

	return sym.(IPlugin), nil
}
