package validation_test

var (
	// schemaDoc is a predefined schema document that covers many different types, useful for
	// testing all parts of this library's SDL-related functionality.
	schemaDoc = []byte(`
		interface Being {
			name(surname: Boolean): String
		}

		interface Pet {
			name(surname: Boolean): String
		}

		interface Canine {
			name(surname: Boolean): String
		}

		enum DogCommand {
			SIT
			HEEL
			DOWN
		}

		type Dog implements Being & Pet & Canine {
			name(surname: Boolean): String
			nickname: String
			barkVolume: Int
			barks: Boolean
			doesKnownCommand(dogCommand: DogCommand): Boolean
			isHousetrained(atOtherHomes: Boolean = true): Boolean
			isAtLocation(x: Int, y: Int): Boolean
		}

		type Cat implements Being & Pet {
			name(surname: Boolean): String
			nickname: String
			meows: Boolean
			meowVolume: Int
			furColor: FurColor
		}

		union CatOrDog = Dog | Cat

		interface Intelligent {
			iq: Int
		}

		type Human implements Being & Intelligent {
			name(surname: Boolean): String
			pets: [Pet]
			relatives: [Human]
			iq: Int
		}

		type Alien implements Being & Intelligent {
			iq: Int
			name(surname: Boolean): String
			numEyes: Int
		}

		union DogOrHuman = Dog | Human

		union HumanOrAlien = Human | Alien

		enum FurColor {
			BROWN
			BLACK
			TAN
			SPOTTED
			NO_FUR
			UNKNOWN
		}

		input ComplexInput {
			requiredField: Boolean!
			nonNullField: Boolean! = false
			intField: Int
			stringField: String
			booleanField: Boolean
			stringListField: [String]
		}

		type ComplicatedArgs {
			intArgField(intArg: Int): String
			nonNullIntArgField(nonNullIntArg: Int!): String
			stringArgField(stringArg: String): String
			booleanArgField(booleanArg: Boolean): String
			enumArgField(enumArg: FurColor): String
			floatArgField(floatArg: Float): String
			idArgField(idArg: ID): String
			stringListArgField(stringListArg: [String]): String
			stringListNonNullArgField(stringListNonNullArg: [String!]): String
			complexArgField(complexArg: ComplexInput): String
			multipleReqs(req1: Int!, req2: Int!): String
			nonNullFieldWithDefault(arg: Int! = 0): String
			multipleOpts(opt1: Int = 0, opt2: Int = 0): String
			multipleOptAndReq(req1: Int!, req2: Int!, opt1: Int = 0, opt2: Int = 0): String
		}

		scalar Invalid

		scalar Any

		type QueryRoot {
			human(id: ID): Human
			alien: Alien
			dog: Dog
			cat: Cat
			pet: Pet
			catOrDog: CatOrDog
			dogOrHuman: DogOrHuman
			humanOrAlien: HumanOrAlien
			complicatedArgs: ComplicatedArgs
			invalidArg(arg: Invalid): String
			anyArg(arg: Any): String
		}

		schema {
			query: QueryRoot
		}

		directive @onQuery on QUERY
		directive @onMutation on MUTATION
		directive @onSubscription on SUBSCRIPTION
		directive @onField on FIELD
		directive @onFragmentDefinition on FRAGMENT_DEFINITION
		directive @onFragmentSpread on FRAGMENT_SPREAD
		directive @onInlineFragment on INLINE_FRAGMENT
	`)

	// queryDoc ...
	queryDoc = []byte(`
		query Foo($humanId: ID) {
			human(id: $humanId) {
				name(surname: false)
				pets {
					...PetFields
				}
			}
		}

		fragment PetFields on Pet {
			name
		}
	`)
)
