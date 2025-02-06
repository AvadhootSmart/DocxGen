"use client";
import React, { useState } from "react";
import { MultiStepLoader } from "../ui/multi-step-loader";
import { IconSquareRoundedX } from "@tabler/icons-react";

const loadingStates = [
    {
        text: "Buying a condo",
    },
    {
        text: "Travelling in a flight",
    },
    {
        text: "Meeting Tyler Durden",
    },
    {
        text: "He makes soap",
    },
    {
        text: "We goto a bar",
    },
    {
        text: "Start a fight",
    },
    {
        text: "We like it",
    },
    {
        text: "Welcome to F**** C***",
    },
];

export function Loader() {
    const [loading, setLoading] = useState(false);
    return (
        <div className="w-full h-[60vh] flex items-center justify-center">
            {/* Core Loader Modal */}
            <MultiStepLoader loadingStates={loadingStates} loading={loading} duration={2000} />
            {loading && (
                <button
                    className="fixed top-4 right-4 text-black dark:text-white z-[120]"
                    onClick={() => setLoading(false)}
                >
                    <IconSquareRoundedX className="h-10 w-10" />
                </button>
            )}
        </div>
    );
}
