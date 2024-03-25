let fibonacci = fn(x) {
    if (x == 0) {
        0
    } else {
        if (x == 1) {
            return 1;
        } else {
            fibonacci(x - 1) + fibonacci(x - 2);
        }
    }
};

puts(fibonacci(1))
puts(fibonacci(2))
puts(fibonacci(3))
puts(fibonacci(4))
puts(fibonacci(5))
