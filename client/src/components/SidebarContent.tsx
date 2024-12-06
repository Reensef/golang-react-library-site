// SidebarContent.tsx
import React from "react";
import {
  Box,
  CloseButton,
  Flex,
  Text,
  useColorModeValue,
  BoxProps,
} from "@chakra-ui/react";
import NavItem from "./NavItem";

import { FiFile, FiUpload, FiArchive } from "react-icons/fi";

interface SidebarProps extends BoxProps {
  onTableChange: (table: "allFiles" | "upload" | "activity") => void;
  onClose: () => void;
}

const SidebarContent = ({ onClose, onTableChange, ...rest }: SidebarProps) => {
  return (
    <Box
      bg={useColorModeValue("white", "gray.900")}
      borderRight="1px"
      borderRightColor={useColorModeValue("gray.200", "gray.700")}
      w={{ base: "full", md: 60 }}
      pos="fixed"
      h="full"
      {...rest}
    >
      <Flex direction="column" h="full">
        <Flex h="20" alignItems="center" mx="8" justifyContent="space-between">
          <Text fontSize="2xl" fontFamily="monospace" fontWeight="bold">
            BooLib
          </Text>
          <CloseButton
            display={{ base: "flex", md: "none" }}
            onClick={onClose}
          />
        </Flex>

        <NavItem icon={FiFile} onClick={() => onTableChange("allFiles")}>
          All Files
        </NavItem>

        <NavItem icon={FiUpload} onClick={() => onTableChange("upload")}>
          Upload
        </NavItem>

        <Box mt="auto" mb="4" width="full">
          <NavItem icon={FiArchive} onClick={() => onTableChange("activity")}>
            User Activity
          </NavItem>
        </Box>
      </Flex>
    </Box>
  );
};

export default SidebarContent;
