/// <reference types="node" />
import React, { useState, useEffect } from "react";
import {
  Box,
  Text,
  VStack,
  Icon,
  useColorModeValue,
  Button,
  Select,
  Flex,
} from "@chakra-ui/react";
import { FiUpload, FiFile } from "react-icons/fi";
import { MdCheckCircle } from "react-icons/md";
import axios from "axios";

interface TagData {
  id: number;
  name: string;
}

const UploadPage = () => {
  const [file, setFile] = useState<File | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [showSuccess, setShowSuccess] = useState(false);
  const [tags, setTags] = useState<TagData[]>([]);
  const [selectedTag, setSelectedTag] = useState<TagData | null>(null);

  const borderColor = useColorModeValue("gray.200", "gray.700");
  const hoverBg = useColorModeValue("gray.50", "gray.600");

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files?.length) {
      setFile(event.target.files[0]);
    }
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!file) {
      return;
    }

    setIsLoading(true);

    try {
      const token = localStorage.getItem("token");
      const formData = new FormData();
      formData.append("file", file);
      if (selectedTag) {
        formData.append("tag_id", selectedTag.id.toString());
      }

      const response = await fetch("/api/v1/files/upload", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        method: "POST",
        body: formData,
      });

      if (!response.ok) {
        throw new Error(`Ошибка при загрузке файла: ${response.status}`);
      }

      setShowSuccess(true);
    } catch (error) {
      console.error(
        (error as Error).message || "Что-то пошло не так при загрузке файла."
      );
    } finally {
      setIsLoading(false);
    }
  };

  const handleDrop = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    const files = event.dataTransfer.files;
    if (files.length > 0) {
      setFile(files[0]);
    }
  };

  const preventDefaults = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    event.stopPropagation();
  };

  useEffect(() => {
    let timer: NodeJS.Timeout;
    if (showSuccess) {
      timer = setTimeout(() => {
        setShowSuccess(false);
        setFile(null);
      }, 2000);
    }
    return () => clearTimeout(timer);
  }, [showSuccess]);

  useEffect(() => {
    const fetchTags = async () => {
      try {
        const token = localStorage.getItem("token");

        const response = await axios.get("/api/v1/files/tags", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        setTags(response.data.data);
      } catch (error) {
        console.error("Error loading tags from the server: ", error);
      }
    };

    fetchTags();
  }, []);

  const textColor = useColorModeValue("gray.600", "gray.400");

  return (
    <Box
      display="flex"
      justifyContent="center"
      alignItems="center"
      height="80vh"
      bg={useColorModeValue("gray.100", "gray.800")}
    >
      <form onSubmit={handleSubmit}>
        <Flex direction="column" alignItems="center">
          <Select
            placeholder="Select tag"
            width="300px"
            onChange={(e) => {
              const selectedTag = tags.find(
                (tag) => tag.name === e.target.value
              );
              setSelectedTag(selectedTag!);
            }}
          >
            {tags.map((tag) => (
              <option value={tag.name}>{tag.name}</option>
            ))}
          </Select>
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
            mt={4}
          >
            {!file && !showSuccess ? (
              <>
                <VStack spacing={4} height="100%" justifyContent="center">
                  <input
                    type="file"
                    name="file"
                    id="fileInput"
                    style={{ display: "none" }}
                    onChange={handleChange}
                  />
                  <label htmlFor="fileInput">
                    <Icon as={FiUpload} boxSize={12} color="blue.500" />
                    <Text fontSize="lg" color={textColor}>
                      Перетащите файлы сюда или кликните для выбора
                    </Text>
                  </label>
                </VStack>
              </>
            ) : showSuccess ? (
              <VStack spacing={4} height="100%" justifyContent="center">
                <Icon as={MdCheckCircle} boxSize={12} color="green.400" />
                <Text>Успех!</Text>
              </VStack>
            ) : (
              <VStack spacing={4} height="100%" justifyContent="center">
                <Icon as={FiFile} boxSize={12} color="green.400" />
                <Text>{file?.name}</Text>
              </VStack>
            )}
          </Box>
          <Button type="submit" mt={4} disabled={isLoading}>
            Загрузить файл
          </Button>
          <Button
            onClick={() => {
              setFile(null);
              setShowSuccess(false);
            }}
            mt={4}
          >
            Сбросить файл
          </Button>
        </Flex>
      </form>
    </Box>
  );
};

export default UploadPage;
