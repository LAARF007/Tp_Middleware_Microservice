"""Analyzer of the technology-watch pipeline.

Reads ARTICLES.ingested, asks an LLM (or the nollm fallback) for a summary, tags and a
category, and publishes the result on ARTICLES.enriched.
"""
