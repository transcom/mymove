# Using `openapi` CLI to compile API specs

## Problem Statement

We have a specification-first development cycle for our APIs. This means that editing our API - adding endpoints,
editing responses, changing functionality - starts in the YAML file that contains the API definition. From that, we use
`go-swagger` to read our specification and generate Go types for use in our backend.

**The good:** With this model, we can focus on the API design without worrying about how to convert that into usable Go
code - `go-swagger` does this for us. Code is neatly organized into separate packages for each API, so they can function
independently.

**The problem:** Our APIs are all concerned with the same data models, so even though they are _technically_
independent, they are highly related. We're defining the same objects over and over again in our YAML specs. All APIs
have a shipment, a move, an orders object, and the list goes on. When we make one change to these objects, we have to
make changes to each and every YAML file.

This means our YAML files quickly get out of sync. We've had to deal with bugs stemming from this disconnect many times.
This is also hugely redundant - there are hundreds of lines that are essentially identical in each API.

We have to do a lot of manual type conversions in the backend to turn the Swagger-generated Go types in our general
model types. These type conversions are also redundant, and they're another place where we can miss changes that add or
modify fields. However, having shared types between APIs would threaten their ability to function independently. (But
should this even be a concern when the services they use on the backend are so interconnected anyway?)

Lastly, we struggle with maintaining the same standards in each API. Some are more resistant to change, and we don't
have a good method for incrementally standardizing those APIs.

## Considered Solutions

1. Write our full API spec in one YAML file and use `go-swagger` to generate types from that spec (status quo).
   - File structure:

     ```text
     mymove/
     ├── swagger/
     │   ├── prime.yaml
     │   ├── support.yaml
     │   ├── ...
     ```

2. Break our spec up into separate files and share definitions between APIs. Use `go-swagger` to generate types from
the split files.
   - File structure:

     ```text
     mymove/
     ├── swagger/
     │   ├── definitions/
     │   │   ├── move.yaml
     │   │   ├── shipment.yaml
     │   │   ├── ...
     │   ├── prime.yaml   <- includes references to move.yaml, shipment.yaml
     │   ├── support.yaml <- includes the same references
     │   ├── ...
     ```

3. Break our spec up into separate files and share definitions between APIs. Use the `openapi` CLI tool to compile the
separate files into one complete YAML file and use `go-swagger` to generate types from the compiled files.
   - File structure:

     ```text
     mymove/
     ├── swagger/
     │   ├── prime.yaml   <- these are generated files, will not be edited
     │   ├── support.yaml
     │   ├── ...
     ├── swagger-def/
     │   ├── definitions/
     │   │   ├── Move.yaml
     │   │   ├── Shipment.yaml
     │   │   ├── ...
     │   ├── prime.yaml   <- includes references to Move.yaml, Shipment.yaml
     │   ├── support.yaml <- includes the same references
     │   ├── ...
     ```

4. Break up and share definitions in a way that prompts `go-swagger` to share types between APIs.
   - I did not find a method that would actually work for this.

## Decision Outcome

### Chosen Alternative: _Use the `openapi` CLI tool to compile shared API definitions (Option 3)_

This looks like the most complicated solution by far. And for the initial implementation, it is. We have already
introduced the `openapi` tool to the project so that we can preview our API documentation, but now we will be dependent
on it for our development process. We will also have to work in a new folder, so all of our engineers will have to
acclimate to the development cycle.

However, the benefits are significant. The `openapi` compiler dictates a structure that is organized and fairly
intuitive, making it easy to create, find, and reference separate definition files. Like option 2, edits to one file can
apply to all of our APIs. Furthermore, the compiler can handle our files as-are, so we can gradually split our
definitions as we move forward.

Unlike option 2, this method won't change the outward behavior of our APIs. External tools like Load Testing, and
eventually the Prime integration, won't need to change the way they consume our content. This was ultimately the
deciding factor because, even though this option _looks_ more complicated, the overall impact of the switch will be
minimal. Load Testing was also completely non-functional with option 2, and I have not yet figured out how to make it
work.

## Pros and Cons of the Alternatives

### Option 1: Use one YAML file for each API (status quo)

- `+` Same development cycle
- `+` All the information is in one place
- `-` Each YAML file is thousands of lines long
- `-` Difficult to keep our definitions in sync
- `-` Difficult to apply and maintain standards

### Option 2: Use split definitions without compiling into a new file

- `+` Same development cycle - no need to update how we generate code and we'll be working in the same folder
- `+` We can structure our sub-folders however we want to
- `-` Third-party tools won't be able to use our APIs the same way. Integrations will
be challenging.
- `-` No defined structure, so we could implement something non-standard or suboptimal
- `-` If you're not careful, the Go types it generates can be strangely different

### Option 3: Use split definitions and compile them into a complete YAML spec

- `+` With a compiled API spec, third-party tools won't have to change how they integrate with us
- `+` The way `go-swagger` generates code will be the same, so our Go types won't change
- `+` Well-defined structure for the shared files so it's easy to navigate
- `+` Makes use of a tool we were already using for documentation purposes
- `-` New development cycle - different folder, new build process
- `-` Looks complicated at first and requires more folders and files
- `-` We'll be relying on a third-party tool to compile our APIs

### Option 4: Use split definitions and find a way to generate them into shared Go types

- `+` Shared Go types could make things easier for us on the backend
- `-` Changing types would require us to update a huge number of files in our backend packages
- `-` Purely hypothetical - I couldn't figure out how to actually do this
