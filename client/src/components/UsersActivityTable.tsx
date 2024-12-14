import { useState } from "react";
import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
  Button,
} from "@chakra-ui/react";
import { ChevronUpIcon, ChevronDownIcon } from "@chakra-ui/icons";

interface ActionData {
  id: number;
  username: string;
  action: string;
  actionDateTime: string;
}

type SortDirection = "asc" | "desc" | null;

const SortButton = ({
  label,
  sortDirection,
  onClick,
}: {
  label: string;
  sortDirection: SortDirection;
  onClick: () => void;
}) => {
  const getSortIcon = () => {
    if (sortDirection === "asc") return <ChevronUpIcon />;
    if (sortDirection === "desc") return <ChevronDownIcon />;
    return undefined; // Иконка отсутствует
  };

  return (
    <Button
      variant="ghost"
      size="sm"
      ml={-3}
      rightIcon={getSortIcon()}
      onClick={onClick}
    >
      {label}
    </Button>
  );
};

const UsersActivityTable = () => {
  const [sortDirection, setSortDirection] = useState({
    name: null,
    dateTime: null,
  });

  // Массив с данными файлов, типизированный через интерфейс FileData
  const [files /*setFiles*/] = useState<ActionData[]>([
    {
      id: 1,
      username: "File1",
      action: "Downloaded",
      actionDateTime: "05.12.2024",
    },
    {
      id: 2,
      username: "File1",
      action: "Downloaded",
      actionDateTime: "05.12.2024",
    },
    {
      id: 3,
      username: "File1",
      action: "Downloaded",
      actionDateTime: "05.12.2024",
    },
  ]);

  const toggleSort = (column: keyof typeof sortDirection) => {
    setSortDirection((prev) => ({
      ...prev,
      [column]:
        prev[column] === "asc"
          ? "desc"
          : prev[column] === "desc"
          ? null
          : "asc",
    }));
  };

  return (
    <TableContainer>
      <Table variant="simple">
        <Thead>
          <Tr>
            <Th width="30%">
              <SortButton
                label="User"
                sortDirection={sortDirection.name}
                onClick={() => toggleSort("name")}
              />
            </Th>
            <Th width="50%">Action</Th>
            <Th width="20%">
              <SortButton
                label="Date & Time"
                sortDirection={sortDirection.dateTime}
                onClick={() => toggleSort("dateTime")}
              />
            </Th>
          </Tr>
        </Thead>
        <Tbody>
          {files.map((actionDate) => (
            <Tr key={actionDate.id}>
              <Td>{actionDate.username}</Td>
              <Td>{actionDate.action}</Td>
              <Td>{actionDate.actionDateTime}</Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default UsersActivityTable;
