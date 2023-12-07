# Automated Factory
The implementation of the factory can be found in factory.go, tests are implemented in factory_test.go and our report is available under report.pdf.

## Example Output
🌱: booted factory with 2 pick-up stations, 2 assembly stations, 2 welding stations, 2 painting stations, 2 drop-off stations, 2 assembly workers, 2 welding workers, 2 painting workers and 2 transport workers\
📨: taskset 1 received.\
📨: taskset 2 received.\
📨: taskset 3 received.\
📝 ➢ 📤: task pickup steel bar arrived at pickup station 0\
📝 ➢ 🚚: taskset 1 arrived at transportation worker 0\
📝 ➢ 📤: task pickup steel wool arrived at pickup station 1\
📝 ➢ 🚚: taskset 2 arrived at transportation worker 1\
🚚 ➢ 📤: transportation worker 1 arrived at pickup station 1\
🕊️ : Facility pickup 1 is free again\
🚚 ➢ 📤: transportation worker 0 arrived at pickup station 0\
🕊️ : Facility pickup 0 is free again\
📝 ➢ 🔨: task weld steel wool arrived at welding station 0\
📝 ➢ 🧑‍: task weld steel wool arrived at welding worker 0\
📝 ➢ 📤: task pickup steel pot arrived at pickup station 1\
📝 ➢ 🧑‍: task weld steel wool arrived at welding worker 1\
📝 ➢ 🔨: task weld steel bar arrived at welding station 1\
🚚 ➢ 🔨: transportation worker 0 arrived at welding station 1\
🚚 ➢ 🔨: transportation worker 1 arrived at welding station 0\
🧑 ‍➢ 🔨: welding worker 1 arrived at welding station 0\
🧑 ‍➢ 🔨: welding worker 0 arrived at welding station 0\
🔨 ➢ ✅: welding task finished\
🕊️ : Facility welding 0 is free again\
📝 ➢ 👷: task assemble steel wool arrived at assembly worker 0\
📝 ➢ 🦾: task assemble steel wool arrived at assembly station 0\
🚚 ➢ 🦾: transportation worker 1 arrived at assembly station 0\
👷 ➢ 🦾: assembly worker 0 arrived at assembly station 0\
🏠: welding worker 1 arrived at control center\
🏠: welding worker 0 arrived at control center\
📝 ➢ 🧑‍: task weld steel bar arrived at welding worker 0\
📝 ➢ 🧑‍: task weld steel bar arrived at welding worker 1\
🧑 ‍➢ 🔨: welding worker 1 arrived at welding station 1\
🧑 ‍➢ 🔨: welding worker 0 arrived at welding station 1\
🦾 ➢ ✅: assembly task finished\
🕊️ : Facility assembly 0 is free again\
📝 ➢ 🧑‍: task paint steel wool in red arrived at painting worker 0\
📝 ➢ 🎨: task paint steel wool in red arrived at painting station 0\
🚚 ➢ 🎨: transportation worker 1 arrived at painting station 0\
🧑 ➢ 🎨: painting worker 0 arrived at painting station 0\
🏠: assembly worker 0 arrived at control center\
🔨 ➢ ✅: welding task finished\
🕊️ : Facility welding 1 is free again\
📝 ➢ 👷: task assemble steel bar arrived at assembly worker 1\
📝 ➢ 🦾: task assemble steel bar arrived at assembly station 1\
🚚 ➢ 🦾: transportation worker 1 arrived at assembly station 1\
👷 ➢ 🦾: assembly worker 0 arrived at assembly station 1\
🏠: welding worker 1 arrived at control center\
🏠: welding worker 0 arrived at control center\
🎨 ➢ ✅: painting task finished\
🕊️ : Facility painting 0 is free again\
📝 ➢ ✈: task dropoff steel wool arrived at dropoff station 0\
🚚 ➢ ✈: transportation worker 1 arrived at dropoff station 0\
✈ ➢ ✅: dropoff task finished\
🕊️ : Facility dropoff 0 is free again\
🦾 ➢ ✅: assembly task finished\
🕊️ : Facility assembly 1 is free again\
✅: taskset 2 was completed :)\
🏠: painting worker 0 arrived at control center\
📝 ➢ 🎨: task paint steel bar in blue arrived at painting station 1\
📝 ➢ 🧑‍: task paint steel bar in blue arrived at painting worker 1\
🏠: transport worker 1 arrived at control center\
📝 ➢ 🚚: taskset 3 arrived at transportation worker 1\
🚚 ➢ 🎨: transportation worker 0 arrived at painting station 1\
🧑 ➢ 🎨: painting worker 1 arrived at painting station 1\
🏠: assembly worker 1 arrived at control center\
🎨 ➢ ✅: painting task finished\
🕊️ : Facility painting 1 is free again\
🚚 ➢ 📤: transportation worker 1 arrived at pickup station 1\
🕊️ : Facility pickup 1 is free again\
📝 ➢ ✈: task dropoff steel bar arrived at dropoff station 1\
📝 ➢ 🧑‍: task weld steel pot arrived at welding worker 0\
📝 ➢ 🔨: task weld steel pot arrived at welding station 0\
📝 ➢ 🧑‍: task weld steel pot arrived at welding worker 1\
🚚 ➢ 🔨: transportation worker 1 arrived at welding station 0\
🧑 ‍➢ 🔨: welding worker 0 arrived at welding station 0\
🧑 ‍➢ 🔨: welding worker 1 arrived at welding station 0\
🚚 ➢ ✈: transportation worker 0 arrived at dropoff station 1\
✈ ➢ ✅: dropoff task finished\
🕊️ : Facility dropoff 1 is free again\
🏠: painting worker 1 arrived at control center\
✅: taskset 1 was completed :)\
🏠: transport worker 0 arrived at control center\
🔨 ➢ ✅: welding task finished\
🕊️ : Facility welding 0 is free again\
📝 ➢ 👷: task assemble steel pot arrived at assembly worker 0\
📝 ➢ 🦾: task assemble steel pot arrived at assembly station 0\
🏠: welding worker 1 arrived at control center\
🚚 ➢ 🦾: transportation worker 1 arrived at assembly station 0\
👷 ➢ 🦾: assembly worker 0 arrived at assembly station 0\
🏠: welding worker 0 arrived at control center\
🦾 ➢ ✅: assembly task finished\
🕊️ : Facility assembly 0 is free again\
📝 ➢ 🧑‍: task paint steel pot in green arrived at painting worker 0\
📝 ➢ 🎨: task paint steel pot in green arrived at painting station 0\
🚚 ➢ 🎨: transportation worker 0 arrived at painting station 0\
🧑 ➢ 🎨: painting worker 1 arrived at painting station 0\
🏠: assembly worker 0 arrived at control center\
🎨 ➢ ✅: painting task finished\
🕊️ : Facility painting 0 is free again\
📝 ➢ ✈: task dropoff steel pot arrived at dropoff station 0\
🏠: painting worker 0 arrived at control center\
🚚 ➢ ✈: transportation worker 1 arrived at dropoff station 0\
✈ ➢ ✅: dropoff task finished\
🕊️ : Facility dropoff 0 is free again\
✅: taskset 3 was completed :)\
🏠: transport worker 1 arrived at control center\