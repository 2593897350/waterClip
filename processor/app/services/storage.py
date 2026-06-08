from pathlib import Path


def processor_output_dir() -> Path:
    output_dir = Path("var/processor")
    output_dir.mkdir(parents=True, exist_ok=True)
    return output_dir


def resolve_input_path(path: str) -> Path:
    candidate = Path(path)
    if candidate.exists():
        return candidate

    repo_candidate = Path.cwd().parent / path
    if repo_candidate.exists():
        return repo_candidate

    return candidate
