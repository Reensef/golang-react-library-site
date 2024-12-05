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

import { LuFile, LuArchive } from "react-icons/lu";

interface SidebarProps extends BoxProps {
  onTableChange: (table: "files" | "activity") => void;
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
            Logo
          </Text>
          <CloseButton
            display={{ base: "flex", md: "none" }}
            onClick={onClose}
          />
        </Flex>

        <NavItem icon={LuFile} onClick={() => onTableChange("files")}>
          All Files
        </NavItem>

        <Box mt="auto" mb="4" width="full">
          <NavItem icon={LuArchive} onClick={() => onTableChange("activity")}>
            User Activity
          </NavItem>
        </Box>
      </Flex>
    </Box>
  );
};

export default SidebarContent;
