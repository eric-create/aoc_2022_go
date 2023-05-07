

with open('hint.txt') as file:
    content: str = file.read()

rounds: list = []

for line in content.splitlines():
    numbers: list = [int(number) for number in line.split(',')]
    rounds.append(sum(numbers))

print(rounds)
print()
for round in rounds:
    print(round)