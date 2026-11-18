import copy
import json
import os
from dataclasses import dataclass, field
from pathlib import Path


@dataclass(frozen=True)
class Settings:
    nats_url: str
    port: int
    # which LLM: "ollama" (local) or "gemini" / "openai" (Gemini, OpenAI-compatible)
    backend: str
    backend_base_url: str
    model: str
    # smallest delay between two LLM calls, to stay under the provider's per-minute quota
    min_interval_s: float
    remote_api_key: str
    tags: list[str] = ()
    categories: list[str] = ()
    # the shape the LLM must answer in
    llm_response_schema: dict = field(default_factory=dict)


def read_config_file(config_dir: Path) -> dict:
    """Reads configuration from a JSON file
    """
    config_path = config_dir / "configuration.json"
    try:
        with config_path.open(encoding="utf-8") as config_file:
            return json.load(config_file)
    except OSError as error:
        raise SystemExit(f"cannot read {config_path}: {error}")
    except json.JSONDecodeError as error:
        raise SystemExit(f"{config_path} is not valid JSON: {error}")


def read_config_secret(environment_name: str) -> str:
    """Reads a secret, from the environment or from a file: NAME, or NAME_FILE. 
    """
    value = os.environ.get(environment_name, "")
    path = os.environ.get(f"{environment_name}_FILE", "")

    if value and path:
        raise SystemExit(
            f"{environment_name} and {environment_name}_FILE are both set: use one or the other"
        )
    if not path and not value:
        raise SystemExit(
            f"neither {environment_name} nor {environment_name}_FILE is defined"
        )
    if not path:
        return value

    try:
        content = Path(path).read_text(encoding="utf-8").strip()
    except OSError as error:
        raise SystemExit(f"cannot read {environment_name}_FILE ({path}): {error}")
    
    return content


def read_config_value(environment_name: str, config: dict, config_key: str, default):
    """Returns the config value according to a priority : 
    env overrides file overrides default. 
    """
    from_environment = os.environ.get(environment_name)
    if from_environment:
        return from_environment

    if config_key in config:
        return config[config_key]

    return default


def load_settings() -> Settings:
    config = read_config_file(Path(os.environ.get("ANALYZER_CONFIG_DIR", ".")))

    # never from configuration.json: it would be committed
    remote_api_key = read_config_secret("REMOTE_API_KEY")

    tags = list(config["tags"])
    categories = list(config["categories"])

    # fill llm response schema with tags and categories
    llm_response_schema = copy.deepcopy(config["llm_response_schema"])
    try:
        llm_response_schema["properties"]["tags"]["items"]["enum"] = tags
        llm_response_schema["properties"]["category"]["enum"] = categories
    except (KeyError, TypeError):
        raise SystemExit(
            "configuration.json: llm_response_schema must have properties.tags.items and properties.category"
        )

    return Settings(
        nats_url=read_config_value("NATS_URL", config, "nats_url", "nats://localhost:4222"),
        port=int(read_config_value("ANALYZER_PORT", config, "port", 8082)),
        backend=read_config_value("ANALYZER_BACKEND", config, "backend",""),
        backend_base_url=read_config_value("ANALYZER_BACKEND_BASE_URL", config, "backend_base_url", "https://generativelanguage.googleapis.com/v1beta/openai/"),
        model=read_config_value("ANALYZER_MODEL", config, "model", "gemini-3.5-flash-lite"),
        min_interval_s=float(read_config_value("ANALYZER_MIN_INTERVAL_S", config, "min_interval_s", 5)),
        remote_api_key=remote_api_key,
        tags=tags,
        categories=categories,
        llm_response_schema=llm_response_schema,
    )


settings = load_settings()
