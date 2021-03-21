import React from 'react'

export default function Button(props) {
    let buttonClasses = [props.fontSize, props.textColor, "duration-200", "group", "flex", "items-center", "rounded-md", "text-opacity-80", "px-4", "py-2", props.bgColor, props.fontWeight]
    let iconClasses = [props.iconSize, props.iconColor, "visible", "duration-200", "mr-2"]

    if (props.bgHoverColor) buttonClasses.push("hover:" + props.bgHoverColor);
    if (props.textHoverColor) buttonClasses.push("hover:" + props.textHoverColor);
    if (props.iconHoverColor) iconClasses.push("group-hover:" + props.iconHoverColor)
    
    return (
        <button onClick={props.onClick} className={buttonClasses.join(" ")}>
            {props.icon
            ? <ion-icon name={props.icon} class={iconClasses.join(" ")}></ion-icon>
            : null
            }
            {props.text}
        </button>
    )
}

Button.defaultProps = {
    textColor: "text-white",
    iconColor: "text-white",
    iconSize: "text-lg",
    fontSize: "text-sm",
    fontWeight: "font-medium"
}