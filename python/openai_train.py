import json
import openai
import os
import pandas as pd
import numpy as np
from openai.embeddings_utils import get_embedding



def get_prompt_from_menu_item(menu_item):
    prompt = ""
    prompt += "Id:"+ str(menu_item["id"]) + "## "

    name = menu_item.get("name") if menu_item.get("name") is not None else {"en":""}
    prompt += "Name:"+ str(name["en"]) + "## "
    
    description = menu_item.get("description") if menu_item.get("description") is not None else {"en":""}
    prompt += "Description:"+ str(description["en"]) + "## "
    
    prompt += "Price: $" + str(float((menu_item["price_info"]["price"]))/100.0) + "## "
    
    allergies = menu_item.get("allergies") if menu_item.get("allergies") is not None else []
    prompt += "Allergies: " + ",".join(allergies) + "## "
    
    prompt += "Contains Alcohol:" + str(menu_item["contains_alcohol"]) 
    return prompt



train_dir = "menus/train_data"
menus = []
for filename in os.listdir(train_dir):
    if filename.endswith(".json"):
        filepath = os.path.join(train_dir, filename)
        with open(filepath, "r") as f:
            json_data = json.load(f)
            menus.append(json_data)


# Set up the OpenAI API credentials
openai.api_key = "" # replace with your actual OpenAI API key

# Encode the preprocessed data using the OpenAI API
menu_items_preprocessed = []
for menu in menus:
    menu_items = menu["menu"]["items"]
    for menu_item in menu_items:
         menu_items_preprocessed.append(get_prompt_from_menu_item(menu_item))

df = pd.DataFrame.from_dict({"menu_item": menu_items_preprocessed})


embedding_model = "text-embedding-ada-002"
embedding_encoding = "cl100k_base"  # this the encoding for text-embedding-ada-002
max_tokens = 300

df["embedding"] = df.menu_item.apply(lambda x: get_embedding(x, engine=embedding_model))
df.to_csv("data/menu_itmes_with_embeddings.csv")
