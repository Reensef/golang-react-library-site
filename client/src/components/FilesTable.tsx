import React, { useState } from "react";
import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
  Button,
  IconButton,
  Badge,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
} from "@chakra-ui/react";
import { ChevronUpIcon, ChevronDownIcon } from "@chakra-ui/icons";
import { FiMoreVertical } from "react-icons/fi";

type SortDirection = "asc" | "desc" | null;

interface FileData {
  id: number;
  name: string;
  tag: string;
  updated: string;
  size: string;
  downloads: number;
}

const SortButton = ({
  label,
  sortDirection,
  onClick,
}: {
  label: string;
  sortDirection: SortDirection;
  onClick: () => void;
}) => {
  return (
    <Button
      variant="ghost"
      size="sm"
      ml={-3}
      rightIcon={
        sortDirection === "asc" ? (
          <ChevronUpIcon />
        ) : sortDirection === "desc" ? (
          <ChevronDownIcon />
        ) : (
          <ChevronDownIcon color="gray.300" />
        )
      }
      onClick={onClick}
    >
      {label}
    </Button>
  );
};

const ActionsMenu = () => {
  return (
    <Menu>
      <MenuButton
        as={IconButton}
        aria-label="Actions"
        icon={<FiMoreVertical />}
        variant="ghost"
      />
      <MenuList>
        <MenuItem>Download</MenuItem>
        <MenuItem>Rename</MenuItem>
        <MenuItem>Delete</MenuItem>
      </MenuList>
    </Menu>
  );
};

const FilesTable = () => {
  const [sortDirection, setSortDirection] = useState({
    name: null,
    updated: null,
    size: null,
    downloads: null,
  });

  // Массив с данными файлов, типизированный через интерфейс FileData
  const [files /*setFiles*/] = useState<FileData[]>([
    {
      id: 1,
      name: "File1",
      tag: "Tag1",
      updated: "05.12.2024",
      size: "12 MB",
      downloads: 15,
    },
    {
      id: 2,
      name: "File2",
      tag: "Tag2",
      updated: "05.12.2024",
      size: "8 MB",
      downloads: 25,
    },
    {
      id: 3,
      name: "File3",
      tag: "Tag3",
      updated: "05.12.2024",
      size: "20 MB",
      downloads: 5,
    },
  ]);

  const toggleSort = (column: keyof typeof sortDirection) => {
    setSortDirection((prev) => ({
      ...prev,
      [column]: prev[column] === "asc" ? "desc" : "asc",
    }));
  };

  return (
    <TableContainer>
      <Table variant="simple">
        <Thead>
          <Tr>
            <Th width="40%">
              <SortButton
                label="Name"
                sortDirection={sortDirection.name}
                onClick={() => toggleSort("name")}
              />
            </Th>

            <Th width="30%">
              <Menu>
                <MenuButton
                  as={Button}
                  variant="ghost"
                  size="sm"
                  ml={-3}
                  rightIcon={<ChevronDownIcon />}
                >
                  Tags
                </MenuButton>
                <MenuList>
                  <MenuItem>Tag 1</MenuItem>
                  <MenuItem>Tag 2</MenuItem>
                  <MenuItem>Tag 3</MenuItem>
                </MenuList>
              </Menu>
            </Th>

            <Th width="10%">
              <SortButton
                label="Last Updated"
                sortDirection={sortDirection.updated}
                onClick={() => toggleSort("updated")}
              />
            </Th>

            <Th width="10%">
              <SortButton
                label="Size"
                sortDirection={sortDirection.size}
                onClick={() => toggleSort("size")}
              />
            </Th>

            <Th width="5%">
              <SortButton
                label="Downloads"
                sortDirection={sortDirection.downloads}
                onClick={() => toggleSort("downloads")}
              />
            </Th>

            <Th width="5%"></Th>
          </Tr>
        </Thead>
        <Tbody>
          {files.map((file) => (
            <Tr key={file.id}>
              <Td>{file.name}</Td>
              <Td>
                <Badge colorScheme="blue">{file.tag}</Badge>
              </Td>
              <Td>{file.updated}</Td>
              <Td>{file.size}</Td>
              <Td>{file.downloads}</Td>
              <Td>
                <ActionsMenu />
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default FilesTable;
