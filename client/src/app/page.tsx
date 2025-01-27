"use client";
import { PlaceholdersAndVanishInput } from "@/components/ui/placeholders-and-vanish-input";
import axios from "axios";
import { useState } from "react";
import Markdown from "react-markdown";
import { motion as m } from "framer-motion";

export default function Home() {
  const [markdown, setMarkdown] = useState("");
  const [url, setUrl] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  const placeholders = [
    "https://github.com/AvadhootSmart/LCDiary",
    "https://github.com/AvadhootSmart/DevDiary",
  ];

  const fetchData = async () => {
    try {
      const response = await axios.post("http://localhost:6969/process-repo", {
        repo_url: url,
      });

      setMarkdown(response.data.data);
      setLoading(false);
      // console.log(response.data.data);
    } catch (error) {
      console.error(error);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUrl(e.target.value);
  };
  const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // console.log("submitted url", url);
    setLoading(true);
    fetchData();
  };

  return (
    <div className="min-h-screen flex flex-col justify-center bg-neutral-800 items-center px-4">
      <h2 className="mb-10 sm:mb-20 text-xl text-center sm:text-5xl text-white ">
        Write code, not documentations
      </h2>
      <div className="w-full">
        <PlaceholdersAndVanishInput
          placeholders={placeholders}
          onChange={handleChange}
          onSubmit={onSubmit}
        />
      </div>
      <div className="flex">
        {loading && <h1 className="text-white text-2xl ">Loading...</h1>}
        {markdown && (
          <m.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.5, ease: "easeInOut" }}
            className="bg-black rounded-md p-4 text-xl text-white"
          >
            <Markdown>{markdown}</Markdown>
          </m.div>
        )}
      </div>
    </div>
  );
}
