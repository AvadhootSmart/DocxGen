"use client";
import { PlaceholdersAndVanishInput } from "@/components/ui/placeholders-and-vanish-input";
import axios from "axios";
import { useState } from "react";
import Markdown from "react-markdown";
import { motion as m } from "framer-motion";
import { DotPattern } from "@/components/ui/dot-pattern";
import { cn } from "@/lib/utils";

export default function Home() {
    const [markdown, setMarkdown] = useState("");
    const [url, setUrl] = useState<string>("");
    const [loading, setLoading] = useState<boolean>(false);

    const placeholders = [
        "Enter a GitHub Repo URL",
        "https://github.com/AvadhootSmart/LCDiary",
        "https://github.com/AvadhootSmart/DevDiary",
    ];

    const fetchData = async () => {
        try {
            const response = await axios.post(
                `${process.env.NEXT_PUBLIC_BACKEND_URL}/process-repo`,
                {
                    repo_url: url,
                },
            );

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
        setLoading(true);
        if (markdown) {
            setMarkdown("");
        }
        fetchData();
    };

    return (
        <>
            <div className="w-full min-h-screen flex flex-col justify-center bg-black font-Poppins items-center lg:px-4 relative overflow-x-hidden">
                <m.div
                    animate={{
                        opacity: [0, 0.5, 1],
                    }}
                    transition={{
                        duration: 2,
                        repeat: Infinity,
                        repeatType: "reverse",
                    }}
                    className="absolute top-0 left-0 w-full h-full"
                >
                    <DotPattern
                        className={cn(
                            "lg:[mask-image:radial-gradient(700px_circle_at_center,white,transparent)]",
                            "sm:[mask-image:radial-gradient(400px_circle_at_center,white,transparent)]",
                            { "opacity-0 transition-all duration-500 ease-in-out": markdown },
                        )}
                    />
                </m.div>
                <div className="fixed  p-9 top-0 left-0 w-fit">
                    <h1 className="text-white text-2xl top-10 left-8">DocxGen</h1>
                </div>
                <h1 className="sm:mb-20 lg:text-5xl text-center sm:text-4xl text-white ">
                    Write code, not documentations
                </h1>
                <m.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 0.5 }}
                    className="w-full sm:px-10"
                >
                    <PlaceholdersAndVanishInput
                        placeholders={placeholders}
                        onChange={handleChange}
                        onSubmit={onSubmit}
                    />
                </m.div>
                {loading && <h1 className="text-white text-2xl mt-10">Loading...</h1>}
            </div>
            {markdown && (
                <m.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: "auto" }}
                    transition={{ duration: 0.5 }}
                    className="flex px-56 mt-10"
                >
                    <m.div className="bg-gray-700 lg:rounded-md min-h-screen sm:rounded-t-2xl lg:p-4 sm:p-2 lg:text-xl text-white prose">
                        <Markdown>{markdown}</Markdown>
                    </m.div>
                </m.div>
            )}
        </>
    );
}
