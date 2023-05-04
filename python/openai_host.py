
import pandas as pd
import numpy as np
from openai.embeddings_utils import get_embedding, cosine_similarity
import openai
openai.api_key = ""
import gradio as gr
import re

# search through the reviews for a specific product

datafile_path = "data/menu_itmes_with_embeddings.csv"

df = pd.read_csv(datafile_path)
df["embedding"] = df.embedding.apply(eval).apply(np.array)
print("Embeddings loaded successfully")


def func_create_json(menu_item):
    j = {}
    key_values = menu_item.split("##")
    for key_value in key_values:
        key = (key_value.split(":")[0]).strip()
        value = (" ".join(key_value.split(":")[1:])).strip()
        if(key == "Description" or key=="Allergies"):
            value = re.sub(r'[^\w\s]+', '', value)
        j[key] = value
    return j

def suggest_menu_items(product_description, n=5, pprint=True):
    product_embedding = get_embedding(
        product_description,
        engine="text-embedding-ada-002"
    )
    df["similarity"] = df.embedding.apply(lambda x: cosine_similarity(x, product_embedding))

    results = (
        df.sort_values("similarity", ascending=False)
        .head(n)
    )[["menu_item"]]

    df.drop(['similarity'],axis=1, inplace=True)

    if pprint:
        for r in results:
            print(r[:200])
            print()
    menu_items = []
    for r in results['menu_item']:
        menu_items.append(func_create_json(r))
    
    return menu_items

iface = gr.Interface(
    share = False,
    fn=suggest_menu_items,
    inputs=gr.inputs.Textbox(label="Query"),
    outputs=gr.outputs.Textbox(label="Suggestions"),
    title="Menu Item Suggestions",
    description="Enter a query to get menu item suggestions.",
    examples=[["vegetarian pasta"], ["spicy chicken burger"], ["chocolate cake"]],
)

iface.launch()
