require 'erb'
require 'fileutils'

kind = ARGV[0]
@fuzzer_name = ARGV[1]

if not ["go", "afl"].include?(kind)
  puts "Please specify fuzzer kind"
  exit 1
end

if File.directory?(@fuzzer_name)
  puts "Fuzzer already exists with this name"
  exit 1
end

FileUtils.copy_entry("template/#{kind}", @fuzzer_name)

for file in ['start', 'environment', 'README.md']
  location = "#{@fuzzer_name}/#{file}"
  erb_location = "#{location}.erb"
  f = File.open(erb_location)
  contents = f.read

  renderer = ERB.new(contents)
  result = renderer.result()

  File.open(location, 'w') { |out| out.write(result) }
  FileUtils.rm(erb_location)
end
