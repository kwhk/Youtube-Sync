import React from 'react'

export default function Button(props) {
    let buttonClasses = ["duration-200", "group", "flex", "items-center", "rounded-md", "text-opacity-80", "px-4", "py-2", props.bgColor]
    let iconClasses = ["visible", "duration-200", "mr-2"]

    if (props.textColor) {
        buttonClasses.push(props.textColor)
    } else {
        buttonClasses.push("text-white")
    }

    if (props.iconColor) {
        iconClasses.push(props.iconColor)
    } else {
        iconClasses.push("text-white")
    }
   
    // adjust font weight
    if (props.fontWeight) {
        if (Array.isArray(props.fontWeight)) {
            buttonClasses.concat(props.fontWeight)
        } else {
            buttonClasses.push(props.fontWeight)
        }
    } else {
        buttonClasses.push("font-medium")
    }
    
    // adjust font size
    if (props.fontSize) {
        if (Array.isArray(props.fontSize)) {
            buttonClasses.concat(props.fontSize)
        } else {
            buttonClasses.push(props.fontSize)
        }
    } else {
        buttonClasses.push("text-sm")
    }

    if (props.bgHoverColor) buttonClasses.push("hover:" + props.bgHoverColor);
    if (props.textHoverColor) buttonClasses.push("hover:" + props.textHoverColor);
    if (props.iconHoverColor) iconClasses.push("group-hover:" + props.iconHoverColor)

    if (props.iconSize) {
        if (Array.isArray(props.iconSize)) {
            iconClasses.concat(props.iconSize)
        } else {
            iconClasses.push(props.iconSize)
        }
    } else {
        iconClasses.push("text-lg")
    }
    
    return (
        <button onClick={props.onClick} className={buttonClasses.join(" ")}>
            <ion-icon name={props.icon} class={iconClasses.join(" ")}></ion-icon>
            {props.text}
        </button>
    )
}