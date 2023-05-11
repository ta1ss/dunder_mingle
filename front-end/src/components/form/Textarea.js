import { forwardRef } from "react"

const Textarea = forwardRef((props, ref) => {
    return (
        <>
            {props.title && 
            <label htmlFor={props.name} className="form-label">
                {props.title}
            </label>}
            <textarea
                type={props.type}
                className={props.className}
                id={props.name}
                ref={ref}
                name={props.name}
                placeholder={props.placeholder}
                onChange={props.onChange}
                autoComplete={props.autoComplete}
                value={props.value}
                rows={props.rows}
            />
            <div className={props.errorDiv}>{props.errorMsg}</div>
        </>
    )
})

export default Textarea