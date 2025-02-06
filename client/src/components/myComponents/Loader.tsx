"use client";
import React from "react";
import { MultiStepLoader as Loader } from "../ui/multi-step-loader";

const loadingStates = [
    {
        text: "Finding your Github Repository",
    },
    {
        text: "Cloning the repository",
    },
    {
        text: "Extracting the necessary source code from the repository",
    },
    {
        text: "Removing jargon data from the source code",
    },
    {
        text: "Analysing the source code",
    },
    {
        text: "Parsing the source code into LLM readable format",
    },
    {
        text: "Requesting AI to generate the API documentation"
    },
];

export function StepLoader({ loading }: { loading: boolean }) {
    // const [loading, setLoading] = useState(false);
    return (
        <div className="w-full h-[60vh] dark flex items-center justify-center">
            {/* Core Loader Modal */}
            <Loader loadingStates={loadingStates} loading={loading} duration={2000} />
        </div>
    );
}
