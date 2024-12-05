import React, { useState } from "react";
import { Box, useDisclosure } from "@chakra-ui/react";
import SidebarContent from "./components/SidebarContent";
import MobileNav from "./components/MobileNav";
import FilesTable from "./components/FilesTable";
import UsersActivityTable from "./components/UsersActivityTable";

import { useColorModeValue } from "@chakra-ui/react";

function App() {
  const { onOpen, onClose } = useDisclosure();
  const [selectedTable, setSelectedTable] = useState<"files" | "activity">(
    "files"
  );

  const handleTableChange = (table: "files" | "activity") => {
    setSelectedTable(table);
  };

  return (
    <Box minH="100vh" bg={useColorModeValue("gray.100", "gray.900")}>
      {/* Sidebar */}
      <SidebarContent
        onTableChange={handleTableChange}
        onClose={onClose}
        display={{ base: "none", md: "block" }}
      />

      {/* Mobile Navigation */}
      <MobileNav onOpen={onOpen} />

      {/* Main Content Area */}
      <Box ml={{ base: 0, md: 60 }} p="4">
        {selectedTable === "files" ? <FilesTable /> : <UsersActivityTable />}
      </Box>
    </Box>
  );
}

export default App;
