import { Card, CardContent } from '@/components/ui/card'

export function DashboardEmptyState() {
    return (
        <div className="flex flex-col items-center justify-center min-h-[400px] p-6">
            <Card className="w-full max-w-lg">
                <CardContent className="flex flex-col items-center text-center pt-10 pb-10 gap-6">
                    <div className="w-14 h-14 rounded-full bg-muted flex items-center justify-center">
                        <svg
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="none"
                            xmlns="http://www.w3.org/2000/svg"
                            className="text-muted-foreground"
                        >
                            <path
                                d="M9 17H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11a2 2 0 0 1 2 2v3"
                                stroke="currentColor"
                                strokeWidth="1.5"
                                strokeLinecap="round"
                            />
                            <path
                                d="M13 21l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2"
                                stroke="currentColor"
                                strokeWidth="1.5"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                            />
                            <rect
                                x="13"
                                y="13"
                                width="8"
                                height="8"
                                rx="1"
                                stroke="currentColor"
                                strokeWidth="1.5"
                            />
                        </svg>
                    </div>

                    <div className="space-y-1">
                        <h2 className="text-lg font-semibold">Welcome to Veylo</h2>
                        <p className="text-sm text-muted-foreground">
                            Get started by adding your first vehicle or creating an inspection.
                        </p>
                    </div>
                </CardContent>
            </Card>
        </div>
    )
}
