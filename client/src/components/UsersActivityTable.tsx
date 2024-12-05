// UserActivityTable.tsx
import React from "react";
import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
} from "@chakra-ui/react";

const UsersActivityTable = () => {
  return (
    <TableContainer>
      <Table variant="simple">
        <Thead>
          <Tr>
            <Th>User</Th>
            <Th>Activity</Th>
            <Th>Timestamp</Th>
          </Tr>
        </Thead>
        <Tbody>
          <Tr>
            <Td>John</Td>
            <Td>Login</Td>
            <Td>2024-12-04 12:30</Td>
          </Tr>
          <Tr>
            <Td>Anna</Td>
            <Td>Upload File</Td>
            <Td>2024-12-04 13:00</Td>
          </Tr>
          <Tr>
            <Td>Maria</Td>
            <Td>Download File</Td>
            <Td>2024-12-04 14:15</Td>
          </Tr>
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default UsersActivityTable;
