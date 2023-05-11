import { forwardRef } from "react"

const Select = forwardRef((props, ref) => {
    return (
        <>
            {props.title && 
            <label htmlFor={props.name} className="form-label">
                {props.title}
            </label>}
            <select
                type={props.type}
                className={props.className}
                id={props.name}
                ref={ref}
                name={props.name}
                placeholder={props.placeholder}
                onChange={props.onChange}
                autoComplete={props.autoComplete}
                value={props.value}
            >
                {props.options.map((option) => {
                    return (
                        <option key={option.id} value={option.value}>
                            {option.id}
                        </option>
                    )
                })}
            </select>
            <div className={props.errorDiv}>{props.errorMsg}</div>
        </>
    )
})

export default Select