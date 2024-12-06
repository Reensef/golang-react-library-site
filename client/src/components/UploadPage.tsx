import React from "react";
import { Box, Text, VStack, Icon, useColorModeValue } from "@chakra-ui/react";
import { FiUpload } from "react-icons/fi";

const UploadPage = () => {
  const borderColor = useColorModeValue("gray.200", "gray.700");
  const hoverBg = useColorModeValue("gray.50", "gray.600");

  const handleDrop = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    const files = event.dataTransfer.files;
    console.log("Files dropped:", files);
    // Handle file upload logic here
  };

  const preventDefaults = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    event.stopPropagation();
  };

  return (
    <Box
      display="flex"
      justifyContent="center"
      alignItems="center"
      height="80vh"
      bg={useColorModeValue("gray.100", "gray.800")}
    >
      <Box
        width="300px"
        height="300px"
        borderWidth={2}
        borderRadius="lg"
        borderColor={borderColor}
        borderStyle="dashed"
        textAlign="center"
        onDrop={handleDrop}
        onDragOver={preventDefaults}
        onDragEnter={preventDefaults}
        onDragLeave={preventDefaults}
        _hover={{ bg: hoverBg }}
        p={4}
      >
        <VStack spacing={4} height="100%" justifyContent="center">
          <Icon as={FiUpload} boxSize={12} color="blue.500" />
          <Text fontSize="lg" color={useColorModeValue("gray.600", "gray.300")}>
            Drag and drop your files here, or click to browse
          </Text>
        </VStack>
      </Box>
    </Box>
  );
};

export default UploadPage;
