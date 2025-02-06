"use client";
import { PlaceholdersAndVanishInput } from "@/components/ui/placeholders-and-vanish-input";
import axios from "axios";
import { useState } from "react";
import Markdown from "react-markdown";
import { motion as m } from "framer-motion";
import { DotPattern } from "@/components/ui/dot-pattern";
import { cn } from "@/lib/utils";
import { LucideClipboard, LucideGithub } from "lucide-react";
import { IconBrandGithub } from "@tabler/icons-react";
import { MultiStepLoader } from "@/components/ui/multi-step-loader";
import { StepLoader } from "@/components/myComponents/Loader";

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
        <div className="bg-black">
            <div className="w-full min-h-screen flex flex-col justify-center bg-black font-Poppins items-center lg:px-4 relative overflow-x-hidden z-10">
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
                <div className="fixed p-10 top-0 left-0 w-full backdrop-blur z-[150]">
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
                {loading && <StepLoader loading={loading} />}
            </div>
            {markdown && (
                <m.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 0.8 }}
                    className="flex lg:px-56 sm:bottom-40 sm:px-10 min-h-screen items-center justify-center w-full relative z-50 lg:bottom-40"
                >
                    <m.div className="bg-gray-700 lg:rounded-md  sm:rounded-2xl lg:p-4 sm:p-2 lg:text-xl text-white prose prose-headings:invert relative mb-10">
                        <LucideClipboard
                            onClick={() => navigator.clipboard.writeText(markdown)}
                            className="absolute z-50 top-6 right-4 text-neutral-400 text-xl cursor-pointer hover:scale-110 transition"
                        />
                        <Markdown className={"-mt-4 mb-10"}>{markdown}</Markdown>
                    </m.div>
                </m.div>
            )}

            <footer className="sm:p-6 lg:p-10 w-full bg-zinc-900 -mt-40 flex justify-between items-center">
                <h1 className="text-white sm:text-xl">DocxGen</h1>
                <IconBrandGithub
                    className="text-white text-2xl cursor-pointer hover:scale-110 transition"
                    onClick={() =>
                        window.open("https://github.com/AvadhootSmart/DocxGen")
                    }
                />
            </footer>
        </div>
    );
}
