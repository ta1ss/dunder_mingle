const Checkbox = (props) => {
    return (
        <div className="form-check p-0 m-0">
            <input
                id={props.name}
                className={props.className}
                type="checkbox"
                value={props.value}
                name={props.name}
                onChange={props.onChange}
                checked={props.checked}
            />
            {props.title && 
            <label htmlFor={props.name} className="form-label m-1">
                {props.title}
            </label>}
        </div>
    )
}

export default Checkbox