package main

import (
	"fmt"
	"invade/env"
	"invade/fly"
	"invade/io/cmdparser"
	"invade/outman"
	"invade/sim"
	"invade/util"
	"io/ioutil"
	//_ "net/http/pprof"
)

/*
func main_profile() {

	// we need a webserver to get the pprof webserver
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	fmt.Println("hello world")
	var wg sync.WaitGroup
	wg.Add(1)
	//go invmain(wg)
	wg.Wait()

}
*/

func main() {

	// VERSION NUMBER
	version := "0.2.3"

	clp := cmdparser.ParseCommandLine()
	if clp.Silent {
		util.InvadeLogger.SetOutput(ioutil.Discard)
	}

	util.InvadeLogger.Println(fmt.Sprintf("Welcome to InvadeGo %s", version))
	// Get command line arguments

	usedseed := util.SetSeed(clp.Seed) // set seed of random number generator

	// Genome
	util.InvadeLogger.Printf("parsing genome definition %s", clp.Genome)
	genome := cmdparser.ParseRegions(clp.Genome)
	if genome == nil {
		panic("Could not obtain valid genome definition")
	} else {
		util.InvadeLogger.Printf("parsed genome definition, will use: %v", genome)
	}

	// Cluster
	util.InvadeLogger.Printf("parsing cluster definition %s", clp.Cluster)
	cluster := cmdparser.ParseRegions(clp.Cluster)
	if cluster == nil {
		util.InvadeLogger.Printf("no piRNA clusters were provided - will not simulate piRNA clusters")
	} else {
		util.InvadeLogger.Printf("parsed piRNA cluster definitions, will use: %v", cluster)
	}

	// Recombination rates
	util.InvadeLogger.Printf("parsing recombination rates %s", clp.RecRate)
	recrate := cmdparser.ParseRecombination(clp.RecRate)
	if recrate == nil {
		util.InvadeLogger.Printf("no recombination rate provided - will not simulate recombination")
	} else {
		util.InvadeLogger.Printf("parsed recombination rate, will use: %v", recrate)
	}

	util.InvadeLogger.Printf("Setting up environment; genome, piRNA cluster, reference regions, trigger sites, paramutable sites and the recombination rate")
	env.SetupEnvironment(genome, cluster, recrate, clp.MinFitness)
	util.InvadeLogger.Print("Setting up jumper")
	env.SetJumper(clp.U, clp.UC)
	util.InvadeLogger.Print("Setting up fitness function")
	fly.SetupFitness(clp.X, clp.T, clp.Noxcluins, clp.MinFitness)
	util.InvadeLogger.Print("Setting up output manager")
	outman.SetupOutputManager(clp.Steps, clp.ReplicateOffset, clp.FileMHP, clp.FileTally, clp.FileSFS, clp.FileDebug, clp.SampleID)

	// Simulate the thing
	util.InvadeLogger.Print("Commencing simulations")
	outman.WriteInfo(clp.ArgString, usedseed, version)
	sim.SimulateInvasions(clp.BasePop, clp.Popsize, clp.Replicates, clp.Generations)
	outman.Done() // let the output manager know the simulations are done
	util.InvadeLogger.Print("Done - thank you for using InvadeGo")

}
