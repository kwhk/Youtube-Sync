import React, { useState } from 'react'

export const withLoading = (WrappedComponent, loadingMessage) => {
    function HOC(props) {
        const [isLoading, setLoading] = useState(true)

        const setLoadingState = isComponentLoading => {
            setLoading(isComponentLoading)
        }
        
        return (
            <div>
                {isLoading && <h1 className="text-white">{loadingMessage}</h1>}
                <WrappedComponent {...props} setLoading={setLoadingState}/>
            </div>
        )
    }

    return HOC
}

export default withLoading