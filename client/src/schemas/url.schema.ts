import { z } from "zod";

const githubRepoRegex =
  /^(https:\/\/github\.com\/[\w-]+\/[\w-]+|git@github\.com:[\w-]+\/[\w-]+\.git)$/;

export const urlSchema = z.object({
  github_repo_url: z
    .string()
    .regex(githubRepoRegex, "Invalid GitHub Repository URL"),
});
