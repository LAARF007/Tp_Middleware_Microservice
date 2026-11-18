"""Entry point: python -m analyzer."""

import asyncio
import logging

from analyzer import config

logger = logging.getLogger("analyzer")


async def main() -> None:
    settings = config.settings
    logger.info(
        "Settings : %s",
        settings,
    )

if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG, format="%(asctime)s %(levelname)s %(name)s: %(message)s")
    asyncio.run(main())
