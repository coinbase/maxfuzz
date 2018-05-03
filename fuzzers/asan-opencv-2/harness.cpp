#include <opencv2/opencv.hpp>
#include <iterator>
#include <string>
#include <iostream>

using namespace cv;
using namespace std;

int main(int argc, char** argv)
{
  std::ifstream file(argv[1]);
  std::vector<char> data;

  file >> std::noskipws;
  std::copy(std::istream_iterator<char>(file),
    std::istream_iterator<char>(),
    std::back_inserter(data));

  Mat matrixJprg = imdecode(Mat(data), 1);
  return 0;
}
