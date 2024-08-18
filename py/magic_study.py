import magic

file = '/bin/ls'

output = magic.from_file(file)
print(output)

if __name__ == '__main__':
    pass
