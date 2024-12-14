import { /*React,*/ useState } from "react";
import { Box, useDisclosure } from "@chakra-ui/react";
import SidebarContent from "./components/SidebarContent";
import MobileNav from "./components/MobileNav";
import FilesTable from "./components/FilesTable";
import UsersActivityTable from "./components/UsersActivityTable";
import UploadPage from "./components/UploadPage";

import { useColorModeValue } from "@chakra-ui/react";

function App() {
  const { onOpen, onClose } = useDisclosure();
  const [selectedTable, setSelectedTable] = useState<
    "allFiles" | "upload" | "activity"
  >("allFiles");

  const handleTableChange = (table: "allFiles" | "upload" | "activity") => {
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
        {selectedTable === "allFiles" ? (
          <FilesTable />
        ) : selectedTable === "upload" ? (
          <UploadPage />
        ) : selectedTable === "activity" ? (
          <UsersActivityTable />
        ) : (
          <Box></Box>
        )}
      </Box>
    </Box>
  );
}

export default App;
