import { useEffect, useState } from "react";
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
import axios from "axios";

interface FileData {
  id: number;
  name: string;
  tag: string;
  updated_at: string;
  size: string;
  downloads: number;
}

type SortDirection = "asc" | "desc" | null;

interface TagData {
  id: number;
  name: string;
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

const formatFileSize = (sizeInKB: number) => {
  if (sizeInKB >= 1024) {
    const sizeInMB = sizeInKB / 1024;
    if (sizeInMB >= 1024) {
      const sizeInGB = sizeInMB / 1024;
      return `${sizeInGB.toFixed(2)} GB`;
    }
    return `${sizeInMB.toFixed(2)} MB`;
  }
  return `${sizeInKB} KB`;
};

const FilesTable = () => {
  const [files, setFiles] = useState<FileData[]>([]);
  const [tags, setTags] = useState<TagData[]>([]); // Состояние для хранения тегов
  const [selectedTagId, setSelectedTagId] = useState<number | null>(null); // Состояние для выбранного тега
  const [sortDirection, setSortDirection] = useState({
    name: null as SortDirection,
    updated_at: null as SortDirection,
    size: null as SortDirection,
    downloads: null as SortDirection,
  });

  const [sortParams, setSortParams] = useState({
    column: "name",
    direction: "asc",
  });

  interface ServerFileData {
    id: number;
    name: string;
    tag: string;
    creator: { username: string };
    updated_at: string;
    size: number;
    downloads: number;
  }

  const fetchFiles = async () => {
    try {
      const params: { [key: string]: string | number } = {
        sort_by: sortParams.column,
        sort_direction: sortParams.direction,
      };
      if (selectedTagId != null) {
        params["tag_id"] = selectedTagId;
      }
      const token = localStorage.getItem("token");
      const response = await axios.get("/api/v1/files", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        params: params,
      });

      const data = response.data.data;

      const formattedData = data.map((item: ServerFileData) => ({
        id: item.id,
        name: item.name,
        tag: item.tag,
        updated_at: new Date(item.updated_at).toLocaleDateString(),
        size: formatFileSize(item.size),
        downloads: item.downloads,
      }));

      setFiles(formattedData);
    } catch (error) {
      console.error("Error loading data from the server: ", error);
    }
  };

  useEffect(() => {
    fetchFiles();
  }, []);

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

  const toggleSort = (column: keyof typeof sortDirection) => {
    let newDirection: SortDirection;

    if (sortDirection[column] === "asc") {
      newDirection = "desc";
    } else {
      newDirection = "asc";
    }

    setSortDirection(() => ({
      name: column === "name" ? newDirection : null,
      updated_at: column === "updated_at" ? newDirection : null,
      size: column === "size" ? newDirection : null,
      downloads: column === "downloads" ? newDirection : null,
    }));

    if (newDirection) {
      setSortParams({
        column: column,
        direction: newDirection,
      });
    } else {
      setSortParams({
        column: "",
        direction: "",
      });
    }
  };

  const resetTagFilter = () => {
    setSelectedTagId(null);
  };

  const handleFileClick = async (id: number) => {
    try {
      const token = localStorage.getItem("token");
      const response = await axios.get(`/api/v1/files/access/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const link = response.data.data;
      window.open(link, "_blank"); // Открываем ссылку в новом окне
    } catch (error) {
      console.error("Error getting file link: ", error);
    }
  };

  const handleFileDownload = async (id: number) => {
    try {
      const token = localStorage.getItem("token");
      const response = await axios.get(`/api/v1/files/access/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
        responseType: "blob", // Для загрузки файла
      });

      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", "file"); // Укажите имя файла
      document.body.appendChild(link);
      link.click();
    } catch (error) {
      console.error("Error downloading file: ", error);
    }
  };

  const handleFileDelete = async (id: number) => {
    try {
      const token = localStorage.getItem("token");
      await axios.delete(`/api/v1/files/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      fetchFiles();
    } catch (error) {
      console.error("Error deleting file: ", error);
    }
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

            <Th width="20%">
              <Menu>
                <MenuButton
                  as={Button}
                  variant="ghost"
                  size="sm"
                  ml={-3}
                  rightIcon={<ChevronDownIcon />}
                >
                  Tags
                  {/* {selectedTagId != null
                    ? tags.find((tag) => tag.id === selectedTagId)?.name
                    : "Tags"} */}
                </MenuButton>
                <MenuList>
                  {tags.map((tag) => (
                    <MenuItem
                      key={tag.id}
                      onClick={() => {
                        if (selectedTagId === tag.id) {
                          resetTagFilter();
                        } else {
                          setSelectedTagId(tag.id);
                        }
                      }}
                      color={selectedTagId === tag.id ? "orange" : "gray.500"}
                    >
                      {tag.name}
                    </MenuItem>
                  ))}
                </MenuList>
              </Menu>
            </Th>

            <Th width="15%">
              <SortButton
                label="Last Updated"
                sortDirection={sortDirection.updated_at}
                onClick={() => toggleSort("updated_at")}
              />
            </Th>

            <Th width="15%">
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
            <Tr key={file.id} _hover={{ bg: "gray.200" }}>
              <Td onClick={() => handleFileClick(file.id)} cursor="pointer">
                {file.name}
              </Td>
              <Td>
                <Badge colorScheme="blue">{file.tag}</Badge>
              </Td>
              <Td>{file.updated_at}</Td>
              <Td>{file.size}</Td>
              <Td>{file.downloads}</Td>
              <Td>
                <Menu>
                  <MenuButton
                    as={IconButton}
                    aria-label="Actions"
                    icon={<FiMoreVertical />}
                    variant="ghost"
                  />
                  <MenuList>
                    <MenuItem onClick={() => handleFileDownload(file.id)}>
                      Download
                    </MenuItem>
                    <MenuItem onClick={() => handleFileDelete(file.id)}>
                      Delete
                    </MenuItem>
                  </MenuList>
                </Menu>
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default FilesTable;
