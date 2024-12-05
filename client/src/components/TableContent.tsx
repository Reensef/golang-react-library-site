// import React from 'react'
import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
} from "@chakra-ui/react";

const TableContent = () => {
  return (
    <TableContainer>
      <Table variant="simple">
        <Thead>
          <Tr>
            <Th>Имя</Th>
            <Th>Возраст</Th>
            <Th>Город</Th>
          </Tr>
        </Thead>
        <Tbody>
          <Tr>
            <Td>Иван</Td>
            <Td>25</Td>
            <Td>Москва</Td>
          </Tr>
          <Tr>
            <Td>Анна</Td>
            <Td>30</Td>
            <Td>Санкт-Петербург</Td>
          </Tr>
          <Tr>
            <Td>Сергей</Td>
            <Td>35</Td>
            <Td>Казань</Td>
          </Tr>
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default TableContent;
