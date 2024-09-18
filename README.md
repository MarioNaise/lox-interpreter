# Lox Interpreter

Check out [Crafting Interpreters](https://craftinginterpreters.com/).

## Lox Syntax

- **Variable Declarations**: Declare variables using the `var` keyword.

  ```lox
  var a = 1;
  ```

- **Control Flow**:

  - `if` and `else`

  ```lox
  if (a > 0) {
    print("Positive");
  } else {
    print("Non-positive");
  }
  ```

  - `for` and `while`

  ```lox
  for (var i = 0; i < 10; i = i + 1) {
      print(i);
  }

  while (a > 0) {
      a = a - 1;
  }
  ```

  - `and` and `or`

  ```lox
  print(true and false);         // false
  print(false and true);         // false
  print(false or true);          // true
  print("thruthy" or true);      // "truthy"
  ```

- **Functions**: Define functions using the `fun` keyword.

  ```lox
  fun greet() {
      print("Hello, World!");
  }
  ```

- **Classes**:

  - Define classes using the `class` keyword.

  ```lox
  class Person() {
      say("Hello!");
  }
  ```

  - Create instances by calling the class.

  ```lox
  var person = Person();

  person.say(); // Hello!
  ```

## Built-in Features

- **Types**: strings, numbers, booleans, and `nil`.

- **Functions**:

  - `read()`: Reads input from the user.
  - `clock()`: Returns the current time in seconds.
  - `print(value)`: Outputs the value to the console.
  - `random(num)`: Creates a random number between 0 and num (num not included).
  - `sleep(milliseconds)`: Pauses execution for the specified duration.
  - `string(value)`: Stringifies the value.
  - `parseNum(string)`: Parses a string to a number.
  - `load(filePath)`: You can load any lox file. Think of loading a file as pasting
    the code directly into the calling file.
    All variables and functions will be available.
    The filepath must be relative to the calling file.

## Commands

- `lox tokenize <filename>`: Prints the scanned tokens from the file.
- `lox parse <filename>`: Prints the parsed AST (kinda improvised).
- `lox evaluate <filename>`: Evaluates a single expression
  from a file (semicolon optional).
- `lox run <filename>` or `lox <filename>`: Runs the file.
- `lox`: Starts an interactive REPL session.
  - Interprets statements and expressions.
    - no multi-line support
    - statements get evaluated
    - expressions get evaluated and printed
  - Enter `.exit` to quit the REPL.

## Getting Started

To run the interpreter, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/MarioNaise/lox-interpreter.git
   cd lox-interpreter
   ```

2. Build the project:

   ```bash
   make
   ```

3. Run the interpreter:

   ```bash
   ./lox
   ```
