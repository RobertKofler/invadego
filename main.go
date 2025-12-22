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
	version := "0.1.8.3"

	clp := cmdparser.ParseCommandLine()
	if clp.Silent {
		util.InvadeLogger.SetOutput(ioutil.Discard)
	}

	util.InvadeLogger.Println(fmt.Sprintf("Welcome to Invade-Spatial %s", version))
	// Get command line arguments

	usedseed := util.SetSeed(clp.Seed) // set seed of random number generator
	// TODO other basis stuff, threads

	// Genome
	util.InvadeLogger.Printf("parsing genome definition %s", clp.Genome)
	genome := cmdparser.ParseRegions(clp.Genome)
	if genome == nil {
		panic("Could not obtain valid genome definition")
	} else {
		util.InvadeLogger.Printf("parsed genome definition, will use: %v", genome)
	}

	// Recombination rates
	util.InvadeLogger.Printf("parsing recombination rates %s", clp.RecRate)
	recrate := cmdparser.ParseRecombination(clp.RecRate)

	// Population Structure
	util.InvadeLogger.Printf("parsing population structure %s", clp.MateBarriers)
	popstruct := cmdparser.ParsePopulationStructureX(clp.MateBarriers, clp.GridX)

	util.InvadeLogger.Printf("Setting up environment; genome, and the recombination rate")
	env.SetupEnvironment(clp.GridX, clp.GridY, clp.ClosedX, clp.ClosedY, genome, recrate, popstruct, clp.BarrierStrength, clp.TriggerSilence, clp.EpiSilencing, clp.PSiLoss, clp.MinFitness, float64(clp.MaxInsertions))
	util.InvadeLogger.Print("Setting up jumper")
	env.SetJumper(clp.U, clp.UC)
	util.InvadeLogger.Print("Setting up fitness function")
	fly.SetupFitness(clp.S, clp.H)
	util.InvadeLogger.Print("Setting up mating function")
	fly.SetupMater(clp.SelfingRate, clp.MateRadius)

	util.InvadeLogger.Print("Setting up output manager")
	om := outman.SetupOutputManager(clp.ConsoleFormat, clp.Steps, clp.ReplicateOffset, clp.FileSpatial, clp.FileMHP, clp.FileDebug, clp.Replicates, clp.SampleID)
	var omint outman.IOutputManager
	if clp.CondInvasion {
		omint = outman.SetupOutputConditionalDecorator(om)
	} else {
		omint = om
	}

	// Simulate the thing
	util.InvadeLogger.Print("Commencing simulations")
	omint.WriteInfo(clp.ArgString, usedseed, version)
	sim.SimulateInvasions(clp.BasePop, clp.Generations, omint)
	omint.End() // let the output manager know the simulations are done
	util.InvadeLogger.Print("Done - thank you for using Invade-Spatial")

}
