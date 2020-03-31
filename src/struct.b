# struct demo

struct Person(
    name, 
    age
);

def main()
    let person = Person("Gabor", 48);
    person.age = person.age + 2;
    print(person);
end
